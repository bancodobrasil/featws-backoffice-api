package services

import (
	"errors"
	"testing"

	"github.com/bancodobrasil/featws-api/dtos"
	"github.com/xanzy/go-gitlab"

	mocks_services "github.com/bancodobrasil/featws-api/mocks/services"
)

func TestMockSave(t *testing.T) {

	fakeRuleSheet := &dtos.Rulesheet{
		ID: 1,
	}

	fakeCommitMessage := "testCommit"

	gitlabMock := new(mocks_services.Gitlab)

	gitlabMock.On("Save", fakeRuleSheet, fakeCommitMessage).Return(nil)

	err := gitlabMock.Save(fakeRuleSheet, fakeCommitMessage)

	if err != nil {
		t.Error("An error occurs on save on gitlab")
	}

}

func TestCreateOrUpdateGilabFileCommitActionFailToDifineFileAction(t *testing.T) {

	// fakeGitlab := mocks_services.NewGitlab(t)

	_, err := CreateOrUpdateGitlabFileCommitAction(&gitlab.Client{}, &gitlab.Project{}, "test", "test", "test")

	gotError := "Failed to define file action:"

	if err.Error() != gotError {
		t.Errorf("Expected %s, got %s", gotError, err)
	}

}

func TestMockFill(t *testing.T) {

	fakeRuleSheet := &dtos.Rulesheet{
		ID: 1,
	}

	service := new(mocks_services.Gitlab)

	service.On("Fill", fakeRuleSheet).Return(nil)

	err := service.Fill(fakeRuleSheet)

	if err != nil {
		t.Error("Error on fill")
	}

}

func TestFillErrorOnConnectTo(t *testing.T) {

	fakeRuleSheet := &dtos.Rulesheet{
		ID: 1,
	}
	service := new(mocks_services.Gitlab)

	service.On("Fill", fakeRuleSheet).Return(errors.New("error on fill"))

	err := service.Fill(fakeRuleSheet)

	if err == nil || err.Error() != "error on fill" {
		t.Error("expected error on fill")
	}

}
