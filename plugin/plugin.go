package plugin

var addonconfig = `{
	"name": "mail",
	"title": "Mail Addon",
	"description": "Send emails based on changes to a store",
	"properties": {
		"mailer": {
			"type": "string",
			"enum": ["mailjet", "mailgun", "postmark", "smtp"]
		},
		"smtp": {
			"type": "object",
			"properties": {
				"server": {
					"type": "string",
					"description": "Your mail server"
				},
				"username": {
					"type": "string",
					"description": "Username for your mail server"
				},
				"password": {
					"type": "string",
					"description": "Password for your mail server"
				},
				"ssl": {
					"type": "boolean",
					"description": "Use an ssl connection for connecting to your mail server"
				}
			},
			"required": ["server", "username", "password"],
			"additionalProperties": false
		},
		"mailgun": {
			"type": "object",
			"properties": {
				"domain": {
					"type": "string"
				},
				"key": {
					"type": "string"
				},
				"public": {
					"type": "string"
				}
			},
			"required": ["domain", "key"],
			"additionalProperties": false
		},
		"postmark": {
			"type": "object",
			"properties": {
				"serverToken": {
					"type": "string"
				},
				"apiToken": {
					"type": "string"
				}
			},
			"required": ["serverToken"],
			"additionalProperties": false
		},
		"mailjet": {
			"type": "object",
			"properties": {
				"apiKey": {
					"type": "string"
				},
				"secretKey": {
					"type": "string"
				}
			},
			"required": ["apiKey"],
			"additionalProperties": false
		}
	},
	"type": "object",
	"additionalProperties": false
}`

// TODO: add jsonschema template format
var linkparams = `{
	"name": "params",
	"title": "Mail Params",
	"properties":{
		"bodyTemplate": {
			"type": "string"
		}, 
		"subjectTemplate": {
			"type": "string",
			"description": "Subject template"
		}, 
		"sender": {
			"type": "string"
		},
		"recipientTemplate":{
			"type": "string"
		}
	},
	"required": [
		"bodyTemplate", "subjectTemplate", "recipientTemplate", "sender"
	],
	"type": "object",
	"additionalProperties": false
}`

// AddonRegistrar an addon registrar
type AddonRegistrar interface {
	Add(name, config, params string)
}

// Register injects an addon into a registry
func Register(ar AddonRegistrar) {
	ar.Add("email", addonconfig, linkparams)
}
