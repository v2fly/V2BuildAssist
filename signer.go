package V2BuildAssist

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/xiaokangwang/VSign/sign"
	"github.com/xiaokangwang/VSign/signerVerify"
	"io"
	"io/ioutil"

	"github.com/xiaokangwang/VSign/insmgr"
)

func RequestForSign(githubToken,
	signerOwner, signerRepo,
	project,
	projectOwner, projectRepo,
	version,
	keyPassword, keyEncrypted string) (int64, string, error) {
	data, id, err := GetReleaseFile(githubToken, projectOwner, projectRepo, version, "Release.unsigned")
	if err != nil {
		return 0, "", err
	}
	//Sort First, this will also read all ins, if some are invalid, it will crash here
	sorted := bytes.NewBuffer(nil)
	insmgr.SortAll(bytes.NewReader(data), sorted)
	//Check version and project
	insall := insmgr.ReadAllIns(bytes.NewReader(sorted.Bytes()))
	if !signerVerify.CheckVersionAndProject(insall, version, project) {
		return 0, "", io.EOF
	}
	r := base64.NewDecoder(base64.StdEncoding, bytes.NewReader([]byte(keyEncrypted)))
	keyorig, err := ioutil.ReadAll(r)
	if err != nil {
		return 0, "", err
	}
	sr, err := sign.Sign(keyorig, keyPassword, sorted.Bytes())
	if err != nil {
		return 0, "", err
	}
	sfpath := fmt.Sprintf("%v/%v.Release", project, version)
	releasebuf := bytes.NewBuffer(nil)
	releasebuf.WriteString("untrusted comment: Signed Release File\n")
	releasebuf.WriteString(string(sr))
	releasebuf.WriteString("\n")
	releasebuf.Write(sorted.Bytes())
	url, err := CreateFileIfNotExist(githubToken, signerOwner, signerRepo, sfpath, "Signed Release", releasebuf.Bytes())
	if err != nil {
		return 0, "", err
	}
	return id, url, nil
}