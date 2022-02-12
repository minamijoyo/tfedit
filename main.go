package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("filename is required.")
	}

	filename := os.Args[1]
	src, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("failed to read file: %s", err)
	}

	f, diags := hclwrite.ParseConfig(src, filename, hcl.Pos{Line: 1, Column: 1})

	if diags.HasErrors() {
		log.Fatalf("failed to parse input: %s", diags)
	}

	blocks := f.Body().Blocks()
	for _, b := range blocks {
		if b.Type() != "resource" {
			continue
		}
		labels := b.Labels()
		if labels[0] != "aws_s3_bucket" {
			continue
		}
		if b.Body().GetAttribute("acl") != nil {
			name := labels[1]
			f.Body().AppendNewline()
			newblock := f.Body().AppendNewBlock("resource", []string{"aws_s3_bucket_acl", name})
			newblock.Body().SetAttributeTraversal("bucket", hcl.Traversal{
				hcl.TraverseRoot{Name: "aws_s3_bucket"},
				hcl.TraverseAttr{Name: name},
				hcl.TraverseAttr{Name: "id"},
			})
			newblock.Body().SetAttributeValue("acl", cty.StringVal("private"))
			b.Body().RemoveAttribute("acl")
		}
	}

	updated := f.BuildTokens(nil).Bytes()
	output := hclwrite.Format(updated)

	fmt.Println(string(output))
}
