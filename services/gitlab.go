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

// Gitlab interface defines methods for saving, filling, and connecting to a Gitlab client.
//
// Property:
//   - Save: A method that takes a pointer to a Rulesheet DTO (Data Transfer Object) and a commit message as input parameters and returns an error. This method is responsible for saving the Rulesheet to Gitlab repository with the provided commit message.
//   - Fill: The method is a function that takes a pointer to a `Rulesheet` DTO and fills it with data from a GitLab repository. It returns an error if there's any issue while filling the `Rulesheet`.
//   - Connect: Connect is a method that returns a pointer to a gitlab.Client and an error. It's used to establish a connection to the GitLab server.
type Gitlab interface {
	Save(rulesheet *dtos.Rulesheet, commitMessage string) error
	Fill(rulesheet *dtos.Rulesheet) error
	Connect() (*gitlab.Client, error)
}

// gitlabService struct holds a pointer to a config.Config object.
//
// Property:
//   - cfg: The `cfg` property is a pointer to a `config.Config` struct, which likely contains configuration settings for a GitLab service.
type gitlabService struct {
	cfg *config.Config
}

// NewGitlab creates a new instance of the Gitlab service using the provided configuration.
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
	rulesheet.Version = fmt.Sprintf("%d", version+1)
	if err != nil {
		log.Errorf("Failed to parse version: %v", err)
		return err
	}
	commitAction, err = createOrUpdateGitlabFileCommitAction(git, proj, cfg.GitlabDefaultBranch, "VERSION", rulesheet.Version+"\n")
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

// createOrUpdateGitlabFileCommitAction this function creates or updates a GitLab file commit action with the specified content.
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

// defineCreateOrUpdateGitlabFileAction this function determines whether to create or update a GitLab file based on whether it already
// exists or not.
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

// Fill is a method in a GitLab service that fills a `Rulesheet` struct with data from GitLab. It first
// checks if a GitLab token is provided, and if not, it returns nil. It then connects to GitLab using the
// provided token and fetches the namespace and project associated with the provided GitLab namespace and
// prefix. It fetches the version, features, parameters, and rules data from the project's default branch
// and populates the corresponding fields in the `Rulesheet` struct.
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

	pipeline, response, err := git.Pipelines.GetLatestPipeline(proj.ID, nil, nil)
	if err != nil {
		log.Errorf("Failed to fetch pipeline: %v", err)
	} else {
		rulesheet.PipelineStatus = pipeline.Status
		rulesheet.WebURL = pipeline.WebURL
		log.Infof("Pipeline Status: %v", response.StatusCode)
	}

	return
}

// Connect this method creates a new GitLab client using the GitLab API token and URL provided in the `gs.cfg`
// configuration object. If the client creation is successful, it returns the GitLab client object,
// otherwise it returns an error.
func (gs *gitlabService) Connect() (*gitlab.Client, error) {
	git, err := gitlab.NewClient(gs.cfg.GitlabToken, gitlab.WithBaseURL(gs.cfg.GitlabURL))

	if err != nil {
		log.Errorf("Failed to create client: %v", err)
		return nil, err
	}
	return git, nil
}

// gitlabLoadJSON loads a JSON file from a GitLab project and decodes it into a given Go struct.
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

// gitlabLoadString that takes in a `gitlab.Client` object, a `gitlab.Project` object, a `ref` string,
// and a `fileName` string as parameters. The function uses the `gitlab` package to fetch a file from a
// GitLab repository using the provided `git` and `proj` objects, with the specified `ref` and `fileName`.
// If the file is not found, an empty byte slice is returned. If the file is found, it is decoded from base64
// and returned as a byte slice.
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
