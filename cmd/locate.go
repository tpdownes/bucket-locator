// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/spf13/cobra"
)

var secretAccessKeyFile string
var secretAccessKeyIdFile string

// locateCmd represents the locate command
var locateCmd = &cobra.Command{
	Use:   "locate",
	Args:  cobra.ExactArgs(1),
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: locateBucket,
}

func init() {
	rootCmd.AddCommand(locateCmd)
	locateCmd.Flags().StringVar(&secretAccessKeyFile, "secret-access-key-file", "", "secret key file")
	locateCmd.Flags().StringVar(&secretAccessKeyIdFile, "secret-access-key-id-file", "", "key id file")
	locateCmd.MarkFlagRequired("secret-access-key-file")
	locateCmd.MarkFlagRequired("secret-access-key-id-file")
}

func locateBucket(cmd *cobra.Command, args []string) {
	bucket := aws.String(args[0])

	secretAccessKeyId, err := ioutil.ReadFile(secretAccessKeyIdFile)
	secretAccessKey, err := ioutil.ReadFile(secretAccessKeyFile)
	svc := s3.New(session.New(), &aws.Config{
		Region:   aws.String("us-central1"),
		Endpoint: aws.String("storage.googleapis.com"),
		Credentials: credentials.NewStaticCredentials(
			strings.TrimSpace(string(secretAccessKeyId)),
			strings.TrimSpace(string(secretAccessKey)),
			""),
	})
	input := &s3.GetBucketLocationInput{
		Bucket: bucket,
	}

	result, err := svc.GetBucketLocation(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}
	fmt.Println(result)
}
