package services

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/bancodobrasil/featws-api/config"
	"github.com/bancodobrasil/featws-api/dtos"
	"github.com/xanzy/go-gitlab"

	"github.com/bancodobrasil/go-featws"
	log "github.com/sirupsen/logrus"
)

// Gitlab ...
type Gitlab interface {
	Save(rulesheet *dtos.Rulesheet, commitMessage string) error
	Fill(rulesheet *dtos.Rulesheet) error
	Connect() (*gitlab.Client, error)
}

type gitlabService struct {
	cfg *config.Config
}

// NewGitlab ...
func NewGitlab(cfg *config.Config) Gitlab {
	return &gitlabService{
		cfg: cfg,
	}
}

func (gs *gitlabService) Save(rulesheet *dtos.Rulesheet, commitMessage string) error {

	cfg := gs.cfg

	if gs.cfg.GitlabToken == "" {
		return nil
	}

	git, err := gs.Connect()
	if err != nil {
		log.Errorf("Error on connect the gitlab client: %v", err)
		return err
	}

	ns, _, err := git.Namespaces.GetNamespace(gs.cfg.GitlabNamespace)
	if err != nil {
		log.Errorf("Failed to fetch namespace: %v", err)
		return err
	}

	proj, resp, err := git.Projects.GetProject(fmt.Sprintf("%s/%s%s", ns.FullPath, gs.cfg.GitlabPrefix, rulesheet.Slug), &gitlab.GetProjectOptions{})
	if err != nil {
		if resp.StatusCode != http.StatusNotFound {
			log.Errorf("Failed to fetch project: %v", err)
			return err
		}

		proj, _, err = git.Projects.CreateProject(&gitlab.CreateProjectOptions{
			Name:        gitlab.String(fmt.Sprintf("%s%s", cfg.GitlabPrefix, rulesheet.Slug)),
			NamespaceID: &ns.ID,
		})
		if err != nil {
			log.Errorf("Failed to create project: %v", err)
			return err
		}
	}

	// projData, _ := json.Marshal(proj)
	// fmt.Println(string(projData))

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
		empty := make([]map[string]interface{}, 0)
		rulesheet.Features = &empty
	}

	sort.Slice(*rulesheet.Features, func(i, j int) bool {
		a := (*rulesheet.Features)[i]
		b := (*rulesheet.Features)[j]
		aValue := a["name"].(string)
		bValue := b["name"].(string)
		return aValue < bValue

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
		empty := make([]map[string]interface{}, 0)
		rulesheet.Parameters = &empty
	}

	sort.Slice(*rulesheet.Parameters, func(i, j int) bool {
		a := (*rulesheet.Parameters)[i]
		b := (*rulesheet.Parameters)[j]
		aValue := a["name"].(string)
		bValue := b["name"].(string)
		return aValue < bValue
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

	// Rules
	if rulesheet.Rules == nil {
		empty := make(map[string]interface{}, 0)
		rulesheet.Rules = &empty
	}

	content, err = json.MarshalIndent(rulesheet.Rules, "", "  ")
	if err != nil {
		log.Errorf("Failed to marshal parameters: %v", err)
		return err
	}
	commitAction, err = createOrUpdateGitlabFileCommitAction(git, proj, cfg.GitlabDefaultBranch, "rules.json", string(content))
	if err != nil {
		log.Errorf("Failed to commit parameters: %v", err)
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

// func printRule(rule interface{}, rulesBuffer *bytes.Buffer, ruleName string, isSliceItem bool) error {
// 	ruleNameTag := "[%s]"
// 	if isSliceItem {
// 		ruleNameTag = "[" + ruleNameTag + "]"
// 	}

// 	switch r := rule.(type) {
// 	case *dtos.Rule:
// 		value, err := json.Marshal(r.Value)
// 		if err != nil {
// 			log.Errorf("Failed marshal rule value: %v", err)
// 			return err
// 		}
// 		fmt.Fprintf(rulesBuffer, ruleNameTag+"\ncondition = %s\nvalue = %s\ntype = object\n\n", ruleName, r.Condition, string(value))
// 	case map[string]interface{}:
// 		fmt.Fprintf(rulesBuffer, ruleNameTag+"\n", ruleName)
// 		keys := make([]string, 0)

// 		for k := range r {
// 			// fmt.Printf("RULE k: %s\n", k)
// 			keys = append(keys, k)
// 		}
// 		sort.Strings(keys)
// 		for _, k := range keys {
// 			v := r[k]
// 			fmt.Fprintf(rulesBuffer, "%s = %s\n", k, v)
// 		}
// 		fmt.Fprintf(rulesBuffer, "\n")
// 	default:
// 		fmt.Fprintf(rulesBuffer, "DEFAULT ENTRY %s = %s\n", ruleName, reflect.TypeOf(rule))
// 	}
// 	return nil
// }

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

func (gs *gitlabService) Fill(rulesheet *dtos.Rulesheet) (err error) {
	if gs.cfg.GitlabToken == "" {
		return nil
	}

	git, err := gs.Connect()
	if err != nil {
		log.Errorf("Error on connect the gitlab client: %v", err)
		return
	}

	ns, _, err := git.Namespaces.GetNamespace(gs.cfg.GitlabNamespace)
	if err != nil {
		log.Errorf("Failed to fetch namespace: %v", err)
		return
	}

	proj, _, err := git.Projects.GetProject(fmt.Sprintf("%s/%s%s", ns.FullPath, gs.cfg.GitlabPrefix, rulesheet.Slug), &gitlab.GetProjectOptions{})
	if err != nil {
		log.Errorf("Failed to fetch project: %v", err)
		return
	}

	bVersion, err := gitlabLoadString(git, proj, gs.cfg.GitlabDefaultBranch, "VERSION")
	if err != nil {
		log.Errorf("Failed to fetch version: %v", err)
		return
	}

	rulesheet.Version = strings.Replace(string(bVersion), "\n", "", -1)

	err = gitlabLoadJSON(git, proj, gs.cfg.GitlabDefaultBranch, "features.json", &rulesheet.Features)
	if err != nil {
		log.Errorf("Failed to fetch features: %v", err)
		return
	}

	err = gitlabLoadJSON(git, proj, gs.cfg.GitlabDefaultBranch, "parameters.json", &rulesheet.Parameters)
	if err != nil {
		log.Errorf("Failed to fetch parameters: %v", err)
		return
	}

	bRulesJSON, err := gitlabLoadString(git, proj, gs.cfg.GitlabDefaultBranch, "rules.json")
	if err != nil {
		log.Errorf("Failed to check rules JSON: %v", err)
		return
	}

	if string(bRulesJSON) != "" {
		err = gitlabLoadJSON(git, proj, gs.cfg.GitlabDefaultBranch, "rules.json", &rulesheet.Rules)
		if err != nil {
			log.Errorf("Failed to fetch parameters: %v", err)
			return
		}
	} else {
		bRules, err := gitlabLoadString(git, proj, gs.cfg.GitlabDefaultBranch, "rules.featws")
		if err != nil {
			log.Errorf("Failed to fetch parameters: %v", err)
			return err
		}

		rulesFile, err := featws.Load(bRules)
		if err != nil {
			log.Errorf("Failed to load featws file: %v", err)
			return err
		}

		rules := make(map[string]interface{})

		s := rulesFile.Section("")

		for _, k := range s.Keys() {
			rules[k.Name()] = k.Value()
		}

		for _, sname := range rulesFile.SectionStrings() {
			if sname == featws.DefaultSection {
				continue
			}
			if sname[:1] == "[" {
				continue
			}
			sec := make(map[string]interface{})
			for _, k := range rulesFile.Section(sname).Keys() {
				sec[k.Name()] = k.Value()
			}
			rules[sname] = sec
		}

		for _, aname := range rulesFile.ArrayStrings() {
			a := rulesFile.Array(aname)

			arr := make([]map[string]interface{}, 0)

			for _, s := range a.Sections() {
				sec := make(map[string]interface{})
				for _, k := range s.Keys() {
					sec[k.Name()] = k.Value()
				}
				arr = append(arr, sec)
			}
			rules[aname[1:len(aname)-1]] = arr
		}

		rulesheet.Rules = &rules
	}

	return
}

func (gs *gitlabService) Connect() (*gitlab.Client, error) {
	git, err := gitlab.NewClient(gs.cfg.GitlabToken, gitlab.WithBaseURL(gs.cfg.GitlabURL))

	if err != nil {
		log.Errorf("Failed to create client: %v", err)
		return nil, err
	}
	return git, nil
}

func gitlabLoadJSON(git *gitlab.Client, proj *gitlab.Project, ref string, fileName string, result interface{}) error {
	rawDecodedText, err := gitlabLoadString(git, proj, ref, fileName)
	if err != nil {
		log.Errorf("Error on load the JSON structure: %v", err)
		return err
	}

	if len(rawDecodedText) > 0 {
		json.Unmarshal(rawDecodedText, result)
	}

	return nil
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

// func ConnectToGitlab() (string, error) {
// 	cfg := config.GetConfig()
// 	if cfg.GitlabToken == "" {
// 		return nil
// 	}

// 	_, err := connectGitlab(cfg)
// 	if err != nil {
// 		log.Errorf("Error on connect the gitlab client: %v", err)
// 		return err
// 	}

// 	return "Ok", nil

// }
