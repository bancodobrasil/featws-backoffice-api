package services

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/bancodobrasil/featws-api/config"
	"github.com/bancodobrasil/featws-api/models"
	"github.com/xanzy/go-gitlab"

	log "github.com/sirupsen/logrus"
)

func saveInGitlab(rulesheet *models.Rulesheet, commitMessage string) error {
	cfg := config.GetConfig()

	git, err := connectGitlab(cfg)
	if err != nil {
		log.Errorf("Error on connect the gitlab client: %v", err)
		return err
	}

	ns, _, err := git.Namespaces.GetNamespace(cfg.GitlabNamespace)
	if err != nil {
		log.Errorf("Failed to fetch namespace: %v", err)
		return err
	}

	proj, resp, err := git.Projects.GetProject(fmt.Sprintf("%s/%s%s", ns.Path, cfg.GitlabPrefix, rulesheet.Name), &gitlab.GetProjectOptions{})
	if err != nil {
		if resp.StatusCode != http.StatusNotFound {
			log.Errorf("Failed to fetch project: %v", err)
			return err
		}

		proj, _, err = git.Projects.CreateProject(&gitlab.CreateProjectOptions{
			Name:        gitlab.String(fmt.Sprintf("%s%s", cfg.GitlabPrefix, rulesheet.Name)),
			NamespaceID: &ns.ID,
		})
		if err != nil {
			log.Errorf("Failed to create project: %v", err)
			return err
		}
	}

	_, resp, err = git.RepositoryFiles.GetFile(proj.ID, "VERSION", &gitlab.GetFileOptions{
		Ref: gitlab.String(cfg.GitlabDefaultBranch),
	})
	if err != nil {
		if resp.StatusCode != http.StatusNotFound {
			log.Errorf("Failed to resolve version: %v", err)
			return err
		}

		rulesheet.Version = "0"

	} else {
		bVersion, err := gitlabLoadString(git, proj, cfg.GitlabDefaultBranch, "VERSION")
		if err != nil {
			log.Errorf("Failed to fetch version: %v", err)
			return err
		}

		rulesheet.Version = strings.Replace(string(bVersion), "\n", "", -1)
	}

	actions := []*gitlab.CommitActionOptions{}
	var commitAction *gitlab.CommitActionOptions
	var content []byte

	// VERSION
	version, err := strconv.Atoi(rulesheet.Version)
	if err != nil {
		log.Errorf("Failed to parse version: %v", err)
		return err
	}
	commitAction, err = createOrUpdateGitlabFileCommitAction(git, proj, cfg.GitlabDefaultBranch, "VERSION", fmt.Sprintf("%d\n", version+1))
	if err != nil {
		log.Errorf("Failed to commit version: %v", err)
		return err
	}
	actions = append(actions, commitAction)

	ci := cfg.GitlabCIScript
	commitAction, err = createOrUpdateGitlabFileCommitAction(git, proj, cfg.GitlabDefaultBranch, ".gitlab-ci.yml", ci)
	if err != nil {
		log.Errorf("Failed to commit ci: %v", err)
		return err
	}
	actions = append(actions, commitAction)

	// FEATURES
	if rulesheet.Features == nil {
		empty := make([]interface{}, 0)
		rulesheet.Features = &empty
	}

	sort.Slice(*rulesheet.Features, func(i, j int) bool {
		a := reflect.ValueOf((*rulesheet.Features)[i])
		b := reflect.ValueOf((*rulesheet.Features)[j])
		aKind := a.Kind()
		bKind := b.Kind()
		if aKind == reflect.Map && bKind == reflect.Map {
			aValue := a.MapIndex(reflect.ValueOf("name")).Interface().(string)
			bValue := b.MapIndex(reflect.ValueOf("name")).Interface().(string)
			return aValue < bValue
		}
		return false
	})

	content, err = json.MarshalIndent(rulesheet.Features, "", "  ")
	if err != nil {
		log.Errorf("Failed to marshal features: %v", err)
		return err
	}
	commitAction, err = createOrUpdateGitlabFileCommitAction(git, proj, cfg.GitlabDefaultBranch, "features.json", string(content))
	if err != nil {
		log.Errorf("Failed to commit features: %v", err)
		return err
	}
	actions = append(actions, commitAction)

	// PARAMETERS
	if rulesheet.Parameters == nil {
		empty := make([]interface{}, 0)
		rulesheet.Parameters = &empty
	}

	sort.Slice(*rulesheet.Parameters, func(i, j int) bool {
		a := reflect.ValueOf((*rulesheet.Parameters)[i])
		b := reflect.ValueOf((*rulesheet.Parameters)[j])
		aKind := a.Kind()
		bKind := b.Kind()
		if aKind == reflect.Map && bKind == reflect.Map {
			aValue := a.MapIndex(reflect.ValueOf("name")).Interface().(string)
			bValue := b.MapIndex(reflect.ValueOf("name")).Interface().(string)
			return aValue < bValue
		}
		return false
	})

	content, err = json.MarshalIndent(rulesheet.Parameters, "", "  ")
	if err != nil {
		log.Errorf("Failed to marshal parameters: %v", err)
		return err
	}
	commitAction, err = createOrUpdateGitlabFileCommitAction(git, proj, cfg.GitlabDefaultBranch, "parameters.json", string(content))
	if err != nil {
		log.Errorf("Failed to commit parameters: %v", err)
		return err
	}
	actions = append(actions, commitAction)

	rulesBuffer := bytes.NewBufferString("")
	// RULES
	if rulesheet.Rules == nil {
		empty := make(map[string]string, 0)
		rulesheet.Rules = &empty
	}

	rules := make([]string, 0)

	for k := range *rulesheet.Rules {
		rules = append(rules, k)
	}

	sort.Strings(rules)

	for _, ruleName := range rules {
		fmt.Fprintf(rulesBuffer, "%s = %s\n", ruleName, (*rulesheet.Rules)[ruleName])
	}

	commitAction, err = createOrUpdateGitlabFileCommitAction(git, proj, cfg.GitlabDefaultBranch, "rules.featws", rulesBuffer.String())
	if err != nil {
		log.Errorf("Failed to commit rules: %v", err)
		return err
	}
	actions = append(actions, commitAction)

	_, _, err = git.Commits.CreateCommit(proj.ID, &gitlab.CreateCommitOptions{
		Branch:        &cfg.GitlabDefaultBranch,
		CommitMessage: gitlab.String(commitMessage),
		Actions:       actions,
	})
	if err != nil {
		log.Errorf("Failed to create commit: %v", err)
		return err
	}

	return err
}

func createOrUpdateGitlabFileCommitAction(git *gitlab.Client, proj *gitlab.Project, ref string, filename string, content string) (*gitlab.CommitActionOptions, error) {
	action, err := defineCreateOrUpdateGitlabFileAction(git, proj, ref, filename)
	if err != nil {
		log.Errorf("Failed to define file action: %v", err)
		return nil, err
	}
	return &gitlab.CommitActionOptions{
		Action:   action,
		FilePath: gitlab.String(filename),
		Content:  gitlab.String(content),
	}, nil
}

func defineCreateOrUpdateGitlabFileAction(git *gitlab.Client, proj *gitlab.Project, ref string, fileName string) (*gitlab.FileActionValue, error) {

	_, resp, err := git.RepositoryFiles.GetFile(proj.ID, fileName, &gitlab.GetFileOptions{
		Ref: gitlab.String(ref),
	})
	if err != nil {
		if resp.StatusCode == http.StatusNotFound {
			return gitlab.FileAction(gitlab.FileCreate), nil
		}

		log.Errorf("Failed to fetch file: %v", err)
		return nil, err
	}

	return gitlab.FileAction(gitlab.FileUpdate), nil
}

func fillWithGitlab(rulesheet *models.Rulesheet) (err error) {
	cfg := config.GetConfig()

	git, err := connectGitlab(cfg)
	if err != nil {
		log.Errorf("Error on connect the gitlab client: %v", err)
		return
	}

	ns, _, err := git.Namespaces.GetNamespace(cfg.GitlabNamespace)
	if err != nil {
		log.Errorf("Failed to fetch namespace: %v", err)
		return
	}

	proj, _, err := git.Projects.GetProject(fmt.Sprintf("%s/%s%s", ns.Path, cfg.GitlabPrefix, rulesheet.Name), &gitlab.GetProjectOptions{})
	if err != nil {
		log.Errorf("Failed to fetch project: %v", err)
		return
	}

	bVersion, err := gitlabLoadString(git, proj, cfg.GitlabDefaultBranch, "VERSION")
	if err != nil {
		log.Errorf("Failed to fetch version: %v", err)
		return
	}

	rulesheet.Version = strings.Replace(string(bVersion), "\n", "", -1)

	rulesheet.Features, err = gitlabLoadJSON(git, proj, cfg.GitlabDefaultBranch, "features.json")
	if err != nil {
		log.Errorf("Failed to fetch features: %v", err)
		return
	}

	rulesheet.Parameters, err = gitlabLoadJSON(git, proj, cfg.GitlabDefaultBranch, "parameters.json")
	if err != nil {
		log.Errorf("Failed to fetch parameters: %v", err)
		return
	}

	bRules, err := gitlabLoadString(git, proj, cfg.GitlabDefaultBranch, "rules.featws")
	if err != nil {
		log.Errorf("Failed to fetch parameters: %v", err)
		return
	}

	rulesArr := strings.Split(string(bRules), "\n")

	rules := make(map[string]string)

	for _, line := range rulesArr {
		if line != "" {
			parts := strings.SplitN(line, "=", 2)
			rules[strings.Trim(parts[0], " ")] = strings.Trim(parts[1], " ")
		}
	}

	rulesheet.Rules = &rules

	return
}

func connectGitlab(cfg *config.Config) (*gitlab.Client, error) {
	git, err := gitlab.NewClient(cfg.GitlabToken, gitlab.WithBaseURL(cfg.GitlabURL))
	if err != nil {
		log.Errorf("Failed to create client: %v", err)
		return nil, err
	}
	return git, nil
}

func gitlabLoadJSON(git *gitlab.Client, proj *gitlab.Project, ref string, fileName string) (*[]interface{}, error) {
	rawDecodedText, err := gitlabLoadString(git, proj, ref, fileName)
	if err != nil {
		log.Errorf("Error on load the JSON structure: %v", err)
		return nil, err
	}

	result := &[]interface{}{}

	if len(rawDecodedText) > 0 {
		json.Unmarshal(rawDecodedText, result)
	}

	return result, nil
}

func gitlabLoadString(git *gitlab.Client, proj *gitlab.Project, ref string, fileName string) ([]byte, error) {
	file, resp, err := git.RepositoryFiles.GetFile(proj.ID, fileName, &gitlab.GetFileOptions{
		Ref: gitlab.String(ref),
	})
	if err != nil {
		if resp.StatusCode == http.StatusNotFound {
			return []byte(""), nil
		}

		log.Errorf("Failed to fetch file: %v", err)
		return nil, err
	}

	rawDecodedText, err := base64.StdEncoding.DecodeString(file.Content)
	if err != nil {
		log.Errorf("Failed to decode base64: %v", err)
		return nil, err
	}
	return rawDecodedText, nil
}
