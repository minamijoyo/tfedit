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
		if len(labels) == 2 && labels[0] != "aws_s3_bucket" {
			continue
		}
		bName := labels[1]

		if b.Body().GetAttribute("acl") != nil {
			f.Body().AppendNewline()
			newblock := f.Body().AppendNewBlock("resource", []string{"aws_s3_bucket_acl", bName})
			newblock.Body().SetAttributeTraversal("bucket", hcl.Traversal{
				hcl.TraverseRoot{Name: "aws_s3_bucket"},
				hcl.TraverseAttr{Name: bName},
				hcl.TraverseAttr{Name: "id"},
			})
			newblock.Body().SetAttributeValue("acl", cty.StringVal("private"))
			b.Body().RemoveAttribute("acl")
		}

		if nested := b.Body().FirstMatchingBlock("logging", []string{}); nested != nil {
			f.Body().AppendNewline()
			newblock := f.Body().AppendNewBlock("resource", []string{"aws_s3_bucket_logging", bName})
			newblock.Body().SetAttributeTraversal("bucket", hcl.Traversal{
				hcl.TraverseRoot{Name: "aws_s3_bucket"},
				hcl.TraverseAttr{Name: bName},
				hcl.TraverseAttr{Name: "id"},
			})
			newblock.Body().AppendUnstructuredTokens(nested.Body().BuildTokens(nil))
			b.Body().RemoveBlock(nested)
		}
	}

	updated := f.BuildTokens(nil).Bytes()
	output := hclwrite.Format(updated)

	fmt.Println(string(output))
}
