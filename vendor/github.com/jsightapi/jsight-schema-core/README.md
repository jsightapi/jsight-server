<div align="center">

<div>  
  &nbsp; 
</div>
	
<a href="https://jsight.io" align="left"><img src="./img/jsight-logo.svg" alt="JSight" width="148px"/></a>

# JSight Schema Go Library

  [![Golang](https://badges.aleen42.com/src/golang.svg)](https://go.dev/)
  [![Telegram support](https://img.shields.io/badge/Support-Telegram-blue)](https://t.me/jsight_support)
  [![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](./CONTRIBUTING.md)
  [![License](https://img.shields.io/github/license/jsightapi/jsight-schema-core?colorB=ff0000)](./LICENSE)
  [![JSight on Facebook](https://img.shields.io/badge/Facebook-1877F2?logo=facebook&logoColor=white)](https://www.facebook.com/jsightapi)
  [![JSight on LinkedIn](https://img.shields.io/badge/LinkedIn-0077B5?logo=linkedin&logoColor=white)](https://www.linkedin.com/company/jsightapi/)
  [![Twitter Follow](https://img.shields.io/twitter/follow/jsightapi.svg?style=social)](https://twitter.com/jsightapi)

  <a href="https://www.alchemistaccelerator.com/portfolio?class=29"><img width="300px" src="./img/alchemist.svg" alt="Alchemist Accelerator"/></a>

<div>
  &nbsp;
</div>

  :star: **Star us on GitHub — it motivates us a lot!** :star:

<div>
  &nbsp;
</div>

</div>

**JSight Schema Go Library** is a library that parses schemas written in [JSight Schema](https://jsight.io/docs/jsight-schema-0-3) language.

**JSight Schema language** — you have never described your data so fast. We mean it. [Compare JSight
Schema with JSON Schema.](#scroll--jsight-schema-language).

JSight Schema language specification on the official website:
https://jsight.io/docs/jsight-schema-0-3.  
JSight Schema language specification on GitHub:
https://github.com/jsightapi/specification/tree/main/versions/JSight%20Schema.

The JSight Schema Go Library is currently used as part of the [JSight
Server](https://github.com/jsightapi/jsight-server) and [JSight Online
Editor](https://github.com/jsightapi/online-editor-frontend) applications.

<div>
  &nbsp;
</div>

<div align="center">
  <a href="https://editor.jsight.io"><img src="./img/try-editor.svg" alt="Try now!" width="300px"/></a>
</div>

## :book: &nbsp; Table of Contents

* [Getting Started](#rocket--getting-started)
* [JSight Schema language](#scroll--jsight-schema-language)
* [Versioning](#bookmark_tabs--versioning)
* [Dependencies](#notebook_with_decorative_cover--dependencies)
* [Tests](#test_tube--tests)
* [Contributing](#sunglasses--contributing)
* [Bugs and Feature Requests](#speech_balloon--bugs-and-feature-requests)
* [Support](#grey_question--support)
* [License](#receipt--license)
* [Resources](#book--resources)
* [Partners](#handshake--partners)
* [Acknowledgments](#trophy--acknowledgments)

<div>
  &nbsp;
</div>

## :rocket: &nbsp; Getting Started

#### Prerequisites

- Go >= 1.18 — [install](https://go.dev/doc/install).

#### Installing

<table>
<tbody>
<tr valign="top">
<td>

Download the JSight Schema Go Library source code:

</td>
<td>

```
git clone https://github.com/jsightapi/jsight-schema-core
```

</td>
</tr>
<tr>
</tr>
<tr valign="top">
<td>

Go to the repository folder:

</td>
<td>

```
cd ./jsight-schema-core/
```

</td>
</tr>
<tr valign="top">
<td>

Install development dependencies:

_(Ensure `$GOPATH/bin` is in the `$PATH`)_

</td>
<td>

```
make deps
```

</td>
</tr>
<tr>
</tr>
<tr valign="top">
<td>

Download all dependencies:

</td>
<td>

```
go mod download
```

</td>
</tr>
<tr>
</tr>
<tr valign="top">
<td>

Run automated tests.

If the tests are successful, then everything is fine, and the library is working.

</td>
<td>

```
go test -cover ./...
```

</td>
</tr>
</tbody>
</table>

&nbsp;

> :warning: **SUPPORT:** If you have any problems while launching the JSight Schema Go Library, do
> not hesitate to contact our support, we respond quickly:  
> Email: support@jsight.io  
> Telegram: https://t.me/jsight_support

&nbsp;

## :scroll: &nbsp; JSight Schema language

The JSight Schema language allows you to describe any data structures with incredible speed and
convenience. You can read more about it in the [JSight Schema Language
Specification](https://jsight.io/docs/jsight-schema-0-3).

The JSight Schema language is actively used by the [JSight
API](https://jsight.io/docs/jsight-api-0-3) language, which is designed to describe API. For more
information about JSight Schema within the JSight API, see the [Quick
Tutorial](https://jsight.io/docs/jsight-api-0-3-quick-tutorial/lesson04).

Mentioned below are examples of the same data schemas described using **JSight Schema** and **JSON
Schema**.

<details open><summary>Example 1. The simplest</summary>

<table align="center">
<thead>
<tr>
  <th width="50%">
    JSight Schema 0.3
  </th>
  <th>
    JSON Schema 2020-12
  </th>
</tr>
</thead>
<tbody>
<tr valign="top">
<td>

```
{
  "id"  : 123, // {min: 1}
  "name": "Tom"
}
```

Pay attention to the main feature of the JSight Schema language. **The basis for a data schema is an
example of valid data.** Additional data requirements are specified in C-like comments. This
approach greatly simplifies the data schema and makes it intuitively clear. Practice shows that such
schema is very simple to create, read and edit.

For details, see the JSight Schema Language Specification, in the [section
“EXAMPLE”](https://jsight.io/docs/jsight-schema-0-3/#example).

<div align="center">

:star: **Star us on GitHub — it motivates us a lot!**

</div>

</td>
<td>

```
{
    "type": "object",
    "required": [
        "id",
        "name"
    ],
    "properties": {
        "id": {
            "type": "integer",
            "minimum": 1
        },
        "name": {
            "type": "string"
        }
    },
    "examples": [{
        "id": 123,
        "name": "Tom"
    }]
}
```
	
</td>
</tr>
</tbody>
</table>

</details>

<details><summary>Example 2. Nested arrays and objects</summary>

<table align="center">
<thead>
<tr>
  <th width="50%">
    JSight Schema 0.3
  </th>
  <th>
    JSON Schema 2020-12
  </th>
</tr>
</thead>
<tbody>
<tr valign="top">
<td>

```
{
  "productId": 1,
  "productName": "An ice sculpture",
  "price": 12.50,    // {min: 0, exclusiveMinimum: true}
  "tags": [          // {minItems: 1, optional: true}
    "cold", 
    "ice"
  ],
  "dimensions": {    // {optional: true}
    "length": 7.0,
    "width": 12.0,
    "height": 9.5
  }
}
```

The JSight Schema language is especially useful when describing nested objects and arrays, which are
very common in real life. You simply give an example of a valid array or object, and add small
clarifying comments.

It is much more complicated in JSON Schema and other languages.

For details, see the JSight Schema Language Specification, sections:

- [EXAMPLE requirements for object
  PROPERTIES](https://jsight.io/docs/jsight-schema-0-3/#example-requirements-for-object-properties),
- [EXAMPLE requirements for array
  ELEMENTS](https://jsight.io/docs/jsight-schema-0-3/#example-requirements-for-array-elements).

<div align="center">

:star: **Star us on GitHub — it motivates us a lot!**

</div>

</td>
<td>

```
{
  "type": "object",
  "properties": {
    "productId": {
      "type": "integer"
    },
    "productName": {
      "type": "string"
    },
    "price": {
      "type": "number",
      "exclusiveMinimum": 0
    },
    "tags": {
      "type": "array",
      "items": {
        "type": "string"
      },
      "minItems": 1
    },
    "dimensions": {
      "type": "object",
      "properties": {
        "length": {
          "type": "number"
        },
        "width": {
          "type": "number"
        },
        "height": {
          "type": "number"
        }
      },
      "required": [ "length", "width", "height" ]
    }
  },
  "required": ["productId","productName","price"],
  "examples": [{
    "productId": 1,
    "productName": "An ice sculpture",
    "price": 12.50,
    "tags": [ "cold", "ice" ],
    "dimensions": {
      "length": 7.0,
      "width": 12.0,
      "height": 9.5
    }
  }]
}
```

This example was created based on the official [Getting Started
instructions](https://json-schema.org/learn/getting-started-step-by-step.html) of the JSON-Schema
standard.

</td>
</tr>
</tbody>
</table>

</details>

<details><summary>Example 3. Property description</summary>

<table align="center">
<thead>
<tr>
  <th width="50%">
    JSight Schema 0.3
  </th>
  <th>
    JSON Schema 2020-12
  </th>
</tr>
</thead>
<tbody>
<tr valign="top">
<td>

```
{ // A product from Acme's catalog.
  "productId": 1,                    // The unique id.
  "productName": "An ice sculpture", // Name of the product.
  
  "price": 12.50, /* {
                       min: 0, 
                       exclusiveMinimum: true
                     } 
                     - The price. */

  "tags": [   /* {minItems: 1, optional: true}    
                 - Tags for the product. */
    "cold", 
    "ice"
  ],
  "dimensions": { // {optional: true} - The dimensions.
    "length": 7.0,
    "width": 12.0,
    "height": 9.5
  }
}
```

Here is the same schema as in the previous example, only property descriptions have been added.

Property descriptions are written in C-like comments. If there are *rules* in the comments, then the
property descriptions are separated by a hyphen.

For details, see the JSight Schema Language Specification, [section “Text notes to
RULES”](https://jsight.io/docs/jsight-schema-0-3#text-notes-to-rules).

<div align="center">

:star: **Star us on GitHub — it motivates us a lot!**

</div>

</td>
<td>

```
{
  "type": "object",
  "description": "A product from Acme's catalog.",
  "properties": {
    "productId": {
      "description": "The unique id.",
      "type": "integer"
    },
    "productName": {
      "description": "Name of the product.",
      "type": "string"
    },
    "price": {
      "description": "The price.",
      "type": "number",
      "exclusiveMinimum": 0
    },
    "tags": {
      "description": "Tags for the product.",
      "type": "array",
      "items": {
        "type": "string"
      },
      "minItems": 1
    },
    "dimensions": {
      "description": "The dimensions.",
      "type": "object",
      "properties": {
        "length": {
          "type": "number"
        },
        "width": {
          "type": "number"
        },
        "height": {
          "type": "number"
        }
      },
      "required": ["length", "width", "height"]
    }
  },
  "required": ["productId","productName","price"],
  "examples": [{
    "productId": 1,
    "productName": "An ice sculpture",
    "price": 12.50,
    "tags": [ "cold", "ice" ],
    "dimensions": {
      "length": 7.0,
      "width": 12.0,
      "height": 9.5
    }
  }]
}
```
	
</td>
</tr>
</tbody>
</table>

</details>

<details><summary>Example 4. Built-in data types</summary>

<table align="center">
<thead>
<tr>
  <th width="50%">
    JSight Schema 0.3
  </th>
  <th>
    JSON Schema 2020-12
  </th>
</tr>
</thead>
<tbody>
<tr valign="top">
<td>

```
{
  "id"        : 123,
  "name"      : "Tom",
  "email"     : "tom@cats.com",              // {type: "email"   }
  "website"   : "http://tom.cats.com",       // {type: "uri"     }
  "birthday"  : "2006-01-02",                // {type: "date"    }
  "registered": "2021-01-02T07:23:12+03:00", // {type: "datetime"}
  "uuid": "550e8400-e29b-41d4-a716-446655440000" // {type: "uuid"}
}
```

JSight Schema supports several additional built-in data types that are required for actual work.

Read more about built-in data types in the JSight Schema Language Specification, section [“RULE
"type"”](https://jsight.io/docs/jsight-schema-0-3#rule-type).

<div align="center">

:star: **Star us on GitHub — it motivates us a lot!**

</div>

</td>
<td>

```
{
  "type": "object",
  "properties": {
    "id": {
      "type": "integer",
      "examples": [123]
    },
    "name": {
      "type": "string",
      "examples": ["Tom"]
    },
    "email": {
      "type": "string",
      "format": "email",
      "examples": ["tom@cats.com"]
    },
    "website": {
      "type": "string",
      "format": "uri",
      "examples": ["http://tom.cats.com"]
    },
    "birthday": {
      "type": "string",
      "format": "date",
      "examples": ["2006-01-02"]
    },
    "registered": {
      "type": "string",
      "format": "date-time",
      "examples": ["2021-01-02T07:23:12+03:00"]
    },
    "uuid": {
      "type": "string",
      "format": "uuid",
      "examples": [
        "550e8400-e29b-41d4-a716-446655440000"
      ]
    }
  },
  "required": [
    "id",
    "name",
    "email",
    "website",
    "birthday",
    "registered",
    "uuid"
  ]
}
```
	
</td>
</tr>
</tbody>
</table>

</details>

<div>
  &nbsp;
</div>

## :bookmark_tabs: &nbsp; Versioning

JSight Schema Go Library releases are versioned according to the [Semantic Versioning
2.0.0](https://semver.org/) standard.

```
{MAJOR version}.{MINOR version}.{PATCH version}
```

Releases are located in the branch `main`, and are tagged with a version number, for example,
`v1.0.0`.

The JSight Schema Go Library release history can be seen here:
https://github.com/jsightapi/jsight-schema-core/releases.

<div>
  &nbsp;
</div>

## :notebook_with_decorative_cover: &nbsp; Dependencies

JSight Schema Go Library dependencies are described in the file [go.mod](./go.mod).

<div>
  &nbsp;
</div>

## :test_tube: &nbsp; Tests

To run automated tests, run the following command in the repository root folder:

```
go test -cover ./...
```

<div>
  &nbsp;
</div>

## :sunglasses: &nbsp; Contributing

Contributing is more than just coding. You can help the project in many ways, and we will be very
happy to accept your contribution to our project.

Details of how you can help the project are described in the [CONTRIBUTING.md](./CONTRIBUTING.md)
document.

<div>
  &nbsp;
</div>

### Contributors

<a href="https://github.com/Emptyfruit"><img src="https://avatars.githubusercontent.com/u/14968783?v=4" width="100" height="100" alt=""/></a>
<a href="https://github.com/add2"><img src="https://avatars.githubusercontent.com/u/3954234?v=4" width="100" height="100" alt=""/></a>
<a href="https://github.com/constantine-malyshev"><img src="https://avatars.githubusercontent.com/u/101567029?v=4" width="100" height="100" alt=""/></a>
<a href="https://github.com/dshemin"><img src="https://avatars.githubusercontent.com/u/11780307?v=4" width="100" height="100" alt=""/></a>
<a href="https://github.com/bakhtinphil"><img src="https://avatars.githubusercontent.com/u/104030166?v=4" width="100" height="100" alt=""/></a>

<div>  
  &nbsp; 
</div>

## :speech_balloon: &nbsp; Bugs and Feature Requests

Do you have a bug report or a feature request? 

Please feel free to add a [new
issue](https://github.com/jsightapi/jsight-schema-core/issues/new) or write to us in support:

- Email: [support@jsight.io](mailto:support@jsight.io)
- Telegram: [@jsight_support](https://t.me/jsight_support)

<div>  
  &nbsp; 
</div>

## :grey_question: &nbsp; Support

If something is unclear to you, please contact support; we try to respond within 24 hours. Moreover,
it is critical for us to understand what is unclear from the first instance.

- Email: [support@jsight.io](mailto:support@jsight.io)
- Telegram: [@jsight_support](https://t.me/jsight_support)

<div>  
  &nbsp; 
</div>

## :receipt: &nbsp; License

This project is licensed under the Apache 2.0 License. See the [LICENSE](./LICENSE) file for more
details.

<div>  
  &nbsp; 
</div>

## :book: &nbsp; Resources

- JSight Official Website: https://jsight.io.

### Documents:

- JSight Schema language specification on the official website:
  https://jsight.io/docs/jsight-schema-0-3.  
- JSight Schema language specification on GitHub:
  https://github.com/jsightapi/specification/tree/main/versions/JSight%20Schema.

### Applications that use the JSight Schema Go Library:

- JSight Online Editor: https://editor.jsight.io.
  - JSight Online Editor repo: https://github.com/jsightapi/online-editor-frontend.
- JSight Server repo: https://github.com/jsightapi/jsight-server.

### Publications:

- JSight blog: https://jsight.io/blog.
- Official Facebook page: https://www.facebook.com/jsightapi.
- Official Twitter: https://twitter.com/jsightapi.
- Official Linkedin: https://www.linkedin.com/company/jsightapi.

### Others:

- All JSight repositories: https://github.com/jsightapi.
- GitHub discussions: https://github.com/jsightapi/jsight-schema-core/discussions.

<div>  
  &nbsp; 
</div>

## :handshake: &nbsp; Partners

- We have successfully completed [class #29 of the Alchemist
  Accelerator](https://www.alchemistaccelerator.com/portfolio?class=29).

<div>  
  &nbsp; 
</div>

## :trophy: &nbsp; Acknowledgments

We sincerely thank all those without whom this project would not have been possible:

- [Alchemist Accelerator](https://www.alchemistaccelerator.com/),
- [Lucas Jones](https://github.com/lucasjones) for his
  [reggen](https://github.com/lucasjones/reggen) library,
- [Dave Collins](https://github.com/davecgh) for his
  [go-spew](https://github.com/davecgh/go-spew) library,
- [Stretchr](https://github.com/stretchr) company for [Testify](https://github.com/stretchr/testify)
  and [Objx](https://github.com/stretchr/objx) libraries,
- [Patrick Mézard](https://github.com/pmezard) for his
  [go-difflib](https://github.com/pmezard/go-difflib) library.

<div align="center">

:star: **Star us on GitHub — it motivates us a lot!** :star:

</div>
