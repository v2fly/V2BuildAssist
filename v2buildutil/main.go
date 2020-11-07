package main

import (
	"fmt"
	"github.com/v2fly/V2BuildAssist"
	"github.com/v2fly/VSign/insmgr"
	"github.com/v2fly/VSign/instimp"
	"io/ioutil"
	"os"
	"strconv"
)

func main() {
	argoffset := 1

	outins := insmgr.NewOutputInsMgr(os.Stdout)

	token := os.Getenv("GITHUB_TOKEN")
	owner := os.Getenv("GITHUB_REPO_OWNER")
	name := os.Getenv("GITHUB_REPO_NAME")

	Sowner := os.Getenv("GITHUB_SREPO_OWNER")
	Sname := os.Getenv("GITHUB_SREPO_NAME")
	Skey := os.Getenv("SIGNING_KEY")

	switch os.Args[0+argoffset] {
	case "gen":
		switch os.Args[1+argoffset] {
		case "sort":
			insmgr.SortAll(os.Stdin, os.Stdout)
		case "version":
			insmgr.NewYieldSingle(instimp.NewVersionIns(os.Args[2+argoffset])).InstructionYield(outins)
		case "project":
			insmgr.NewYieldSingle(instimp.NewProjectIns(os.Args[2+argoffset])).InstructionYield(outins)
		case "file":
			instimp.NewFileBasedInsYield(os.Args[2+argoffset]).InstructionYield(outins)
			return
		}
	case "post":
		data, _ := ioutil.ReadAll(os.Stdin)
		switch os.Args[1+argoffset] {
		case "commit":
			fmt.Println(V2BuildAssist.CreateCommentForCommit(token, owner, name, os.Args[2+argoffset], string(data)))
		case "pr":
			i, err := strconv.Atoi(os.Args[2+argoffset])
			if err != nil {
				panic(err)
			}
			fmt.Println(V2BuildAssist.CreateCommentForPR(token, owner, name, string(data), i))
		}
		return
	case "sign":
		password := os.Args[1+argoffset]
		version := os.Args[2+argoffset]
		project := os.Args[3+argoffset]
		fmt.Println(V2BuildAssist.RequestForSign(token, Sowner, Sname, project, owner, name, version, password, Skey))
	}
}
