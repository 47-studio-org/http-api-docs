package docs

import (
	"bytes"
	"fmt"
	"html"
	"regexp"
	"strings"
	"time"
)

// MarkdownFormatter implements a markdown doc generator. It is
// used to generate the IPFS website API reference at
// https://github.com/ipfs/website/blob/master/content/pages/docs/api.md
type MarkdownFormatter struct{}

func (md *MarkdownFormatter) GenerateIntro() string {
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, `---
title: HTTP API
legacyUrl: https://docs.ipfs.io/reference/api/http/
description: HTTP API reference for IPFS, the InterPlanetary File System.
---

# HTTP API reference

<!-- TODO: Describe how to change ports and configure the API server -->
<!-- TODO: Structure this around command groups (dag, object, files, etc.) -->

_Generated on %s, from go-ipfs v%s._

When an IPFS node is running as a daemon, it exposes an HTTP API that allows you to control the node and run the same commands you can from the command line.

In many cases, using this API this is preferable to embedding IPFS directly in your program — it allows you to maintain peer connections that are longer lived than your app and you can keep a single IPFS node running instead of several if your app can be launched multiple times. In fact, the `+"`ipfs`"+` CLI commands use this API when operating in online mode.

::: tip
This document was autogenerated from go-ipfs. For issues and support, check out the [http-api-docs](https://github.com/ipfs/http-api-docs) repository on GitHub.
:::

## Getting started

### Alignment with CLI commands

The HTTP API under `+"`/api/v0/`"+` is an RPC-style API over HTTP, not a REST API.

[Every command](/reference/cli/) usable from the CLI is also available through the HTTP API. For example:
`+"```sh"+
		`
> ipfs swarm peers
/ip4/104.131.131.82/tcp/4001/p2p/QmaCpDMGvV2BGHeYERUEnRQAwe3N8SzbUtfsmvsqQLuvuJ
/ip4/104.236.151.122/tcp/4001/p2p/QmSoLju6m7xTh3DuokvT3886QRYqxAzb1kShaanJgW36yx
/ip4/104.236.176.52/tcp/4001/p2p/QmSoLnSGccFuZQJzRadHn95W2CrSFmZuTdDWP8HXaHca9z

> curl -X POST http://127.0.0.1:5001/api/v0/swarm/peers
{
  "Strings": [
    "/ip4/104.131.131.82/tcp/4001/p2p/QmaCpDMGvV2BGHeYERUEnRQAwe3N8SzbUtfsmvsqQLuvuJ",
    "/ip4/104.236.151.122/tcp/4001/p2p/QmSoLju6m7xTh3DuokvT3886QRYqxAzb1kShaanJgW36yx",
    "/ip4/104.236.176.52/tcp/4001/p2p/QmSoLnSGccFuZQJzRadHn95W2CrSFmZuTdDWP8HXaHca9z",
  ]
}
`+"```"+`

### Arguments

Arguments are added through the special query string key "arg":

`+"```"+`
> curl -X POST "http://127.0.0.1:5001/api/v0/swarm/disconnect?arg=/ip4/54.93.113.247/tcp/48131/p2p/QmUDS3nsBD1X4XK5Jo836fed7SErTyTuQzRqWaiQAyBYMP"
{
  "Strings": [
    "disconnect QmUDS3nsBD1X4XK5Jo836fed7SErTyTuQzRqWaiQAyBYMP success",
  ]
}
`+"```"+`

Note that it can be used multiple times to signify multiple arguments.

### Flags

Flags are added through the query string. For example, the %s flag is the %s query parameter below:

`+"```"+`
> curl -X POST "http://127.0.0.1:5001/api/v0/object/get?arg=QmaaqrHyAQm7gALkRW8DcfGX3u8q9rWKnxEMmf7m9z515w&encoding=json"
{
  "Links": [
    {
      "Name": "index.html",
      "Hash": "QmYftndCvcEiuSZRX7njywX2AGSeHY2ASa7VryCq1mKwEw",
      "Size": 1700
    },
    {
      "Name": "static",
      "Hash": "QmdtWFiasJeh2ymW3TD2cLHYxn1ryTuWoNpwieFyJriGTS",
      "Size": 2428803
    }
  ],
  "Data": "CAE="
}
`+"```"+`

::: tip
Some arguments may belong only to the CLI but appear here too. These usually belong to client-side processing of input, particularly in the `+"`add`"+` command.
:::


## HTTP status codes

Status codes used at the RPC layer are simple:

- `+"`200`"+` - The request was processed or is being processed (streaming)
- `+"`500`"+` - RPC endpoint returned an error
- `+"`400`"+` - Malformed RPC, argument type error, etc
- `+"`403`"+` - RPC call forbidden
- `+"`404`"+` - RPC endpoint doesn't exist
- `+"`405`"+` - HTTP Method Not Allowed

Status code `+"`500`"+` means that the function _does_ exist, but IPFS was not able to fulfil the request because of an error. To know that reason, you have to look at the the error message that is usually returned with the body of the response (if no error, check the daemon logs).

Streaming endpoints fail as above, unless they have started streaming. That means they will have sent a `+"`200`"+` status code already. If an error happens during the stream, it will be included in a Trailer response header (some endpoints may additionally include an error in the last streamed object).

A `+"`405`"+`error may mean that you are using the wrong HTTP method (i.e. GET instead of POST), or that you are not allowed to call that method (i.e. due to CORS restrictions when making a request from a browser).

## HTTP commands
`,
		time.Now().Format("2006-01-02"),
		IPFSVersion(),
		"`--encoding=json`",
		"`&encoding=json`")

	return buf.String()
}

func (md *MarkdownFormatter) GenerateIndex(endps []*Endpoint) string {
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "## Index\n\n")

	for _, endp := range endps {
		fmt.Fprintf(buf, "  *  [%s](#%s)\n",
			strings.TrimPrefix(endp.Name, "/api/v0"),
			strings.Replace(strings.TrimPrefix(endp.Name, "/"), "/", "-", -1))
	}

	buf.WriteString("\n\n## Endpoints\n\n")
	return buf.String()
}

func (md *MarkdownFormatter) GenerateEndpointBlock(endp *Endpoint) string {
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, `
## %s

%s


`, endp.Name, html.EscapeString(endp.Description))
	return buf.String()
}

func (md *MarkdownFormatter) GenerateArgumentsBlock(args []*Argument, opts []*Argument) string {
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "### Arguments\n\n")

	if len(args)+len(opts) == 0 {
		fmt.Fprintf(buf, "This endpoint takes no arguments.\n")
	}

	for _, arg := range args {
		fmt.Fprintf(buf, genArgument(arg, true))
	}
	for _, opt := range opts {
		fmt.Fprintf(buf, genArgument(opt, false))
	}

	fmt.Fprintf(buf, "\n")
	return buf.String()
}

// Removes the "Default:..." part in the descriptions.
var fixDesc, _ = regexp.Compile(" Default: [a-zA-z0-9-_]+ ?\\.")

func genArgument(arg *Argument, aliasToArg bool) string {
	// These get handled by GenerateBodyBlock
	if arg.Type == "file" {
		return "\n"
	}

	buf := new(bytes.Buffer)
	alias := arg.Name
	if aliasToArg {
		alias = "arg"
	}

	fixedDescription := string(fixDesc.ReplaceAll([]byte(arg.Description), []byte("")))
	fixedDescription = html.EscapeString(fixedDescription)

	fmt.Fprintf(buf, "- `%s` [%s]: %s", alias, arg.Type, fixedDescription)
	if len(arg.Default) > 0 {
		fmt.Fprintf(buf, " Default: `%s`.", arg.Default)
	}
	if arg.Required {
		fmt.Fprintf(buf, ` Required: **yes**.`)
	} else {
		fmt.Fprintf(buf, ` Required: no.`)
	}
	fmt.Fprintln(buf)
	return buf.String()
}

func (md *MarkdownFormatter) GenerateBodyBlock(args []*Argument) string {
	var bodyArg *Argument
	for _, arg := range args {
		if arg.Type == "file" {
			bodyArg = arg
			break
		}
	}

	if bodyArg != nil {
		buf := new(bytes.Buffer)
		fmt.Fprintf(buf, `
### Request Body

Argument `+"`%s`"+` is of file type. This endpoint expects one or several files (depending on the command) in the body of the request as 'multipart/form-data'.

`, bodyArg.Name)

		// Special documentation for /add
		if bodyArg.Endpoint == "/api/v0/add" {
			fmt.Fprintln(buf, `

The `+"`add`"+` command not only allows adding files, but also uploading directories and complex hierarchies.

This happens as follows: Every part in the multipart request is a *directory* or a *file* to be added to IPFS.

Directory parts have a special content type `+"`application/x-directory`"+`. These parts do not carry any data. The part headers look as follows:

`+"```"+`
Content-Disposition: form-data; name="file"; filename="folderName"
Content-Type: application/x-directory
`+"```"+`

File parts carry the file payload after the following headers:

`+"```"+`
Abspath: /absolute/path/to/file.txt
Content-Disposition: form-data; name="file"; filename="folderName%2Ffile.txt"
Content-Type: application/octet-stream

...contents...
`+"```"+`

The above file includes its path in the "folderName/file.txt" hierarchy and IPFS will therefore be able to add it inside "folderName". The parts declaring the directories are optional when they have files inside and will be inferred from the filenames. In any case, a depth-first traversal of the directory tree is recommended to order the different parts making the request.

The `+"`Abspath`"+` header is included for filestore/urlstore features that are enabled with the `+"`nocopy`"+` option and it can be set to the location of the file in the filesystem (within the IPFS root), or to its full web URL.
`)
		}
		return buf.String()
	}
	return ""
}

func (md *MarkdownFormatter) GenerateResponseBlock(response string) string {
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, `
### Response

On success, the call to this endpoint will return with 200 and the following body:

`)

	buf.WriteString("```json\n" + response + "\n```\n\n")

	return buf.String()
}

func (md *MarkdownFormatter) GenerateExampleBlock(endp *Endpoint) string {
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "### cURL Example\n\n")
	fmt.Fprintf(buf, "`")
	fmt.Fprintf(buf, "curl -X POST ")

	// Assemble arguments which are not of type file
	var queryargs []string
	hasFileArg := false
	for _, arg := range endp.Arguments {
		q := "arg="
		if arg.Type != "file" {
			q += "<" + arg.Name + ">"
			queryargs = append(queryargs, q)
		} else {
			hasFileArg = true
		}
	}

	// Assemble options
	for _, opt := range endp.Options {
		q := opt.Name + "="
		//if !opt.Required { // Omit non required options
		//	continue
		//}
		if len(opt.Default) > 0 {
			q += opt.Default
		} else {
			q += "<value>"
		}
		queryargs = append(queryargs, q)
	}

	if hasFileArg {
		fmt.Fprintf(buf, "-F file=@myfile ")
	}

	fmt.Fprintf(buf, "\"http://127.0.0.1:5001%s", endp.Name)
	if len(queryargs) > 0 {
		fmt.Fprintf(buf, "?%s\"", strings.Join(queryargs, "&"))
	} else {
		fmt.Fprintf(buf, "\"")
	}

	fmt.Fprintf(buf, "`\n\n---\n")
	return buf.String()
}
