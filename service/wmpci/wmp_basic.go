package wmpci

var version = "v0.0.1"

func Version() string {
	return version
}

type WMPBasic struct {
	Id           string   `json:"id" yaml:"id" bson:"id"`
	Name         string   `json:"name" yaml:"name" bson:"name"`                                                          // plugin full name
	Version      string   `json:"version" yaml:"version" bson:"version"`                                                 // V-0.0.1-1
	PushState    string   `json:"push_state,omitempty" yaml:"push_state,omitempty" bson:"push_state,omitempty"`          // Develop/Abandon/Release
	OpenSource   string   `json:"open_source,omitempty" yaml:"open_source,omitempty" bson:"open_source,omitempty"`       // Describe the open source address of the plug-in
	Description  string   `json:"description,omitempty" yaml:"description,omitempty" bson:"description,omitempty"`       // Plug-in description information
	Certificate  string   `json:"certificate,omitempty" yaml:"certificate,omitempty" bson:"certificate,omitempty"`       // The certificate issued to this plugin
	PluginAvatar string   `json:"plugin_avatar,omitempty" yaml:"plugin_avatar,omitempty" bson:"plugin_avatar,omitempty"` // Plug-in avatar is converted to binary array storage
	Copyright    []Author `json:"copyright,omitempty" yaml:"copyright,omitempty" bson:"copyright,omitempty"`
}

type Author struct {
	AuthorName   string   `json:"author_name" yaml:"author_name" bson:"author_name"`                                     // Describe full name of the author
	AuthorEmail  string   `json:"author_email,omitempty" yaml:"author_email,omitempty" bson:"author_email,omitempty"`    // Describe mail address of the author
	AuthorAvatar string   `json:"author_avatar,omitempty" yaml:"author_avatar,omitempty" bson:"author_avatar,omitempty"` // Author avatar is converted to binary array storage
	Sponsor      []string `json:"sponsor,omitempty" yaml:"sponsor,omitempty" bson:"sponsor,omitempty"`                   // Plug-in sponsorship information list, sponsor through these links
}

func (basic WMPBasic) RegisterBasicInfo(id, name, version, desc string) WMPBasic {
	return WMPBasic{Id: id, Name: name, Version: version, Description: desc}
}

func (basic WMPBasic) RegisterCopyrightInfo(name, email string, avatar string, sponsor []string) WMPBasic {
	basic.Copyright = append(basic.Copyright, Author{
		AuthorName:   name,
		AuthorEmail:  email,
		AuthorAvatar: avatar,
		Sponsor:      sponsor,
	})
	return basic
}

type WMPRequest map[string][]interface{}

type WMPResponse map[string][]interface{}

type WMPCustom struct {
	Name         string      `json:"name" yaml:"name" bson:"name"`                                                          // Name Definitions for Custom Field Attributes
	Value        interface{} `json:"value" yaml:"value" bson:"value"`                                                       // The actual value of the custom field attribute
	Description  string      `json:"description,omitempty" yaml:"description,omitempty" bson:"description,omitempty"`       // Used to describe custom field properties
	FormatRegexp string      `json:"format_regexp,omitempty" yaml:"format_regexp,omitempty" bson:"format_regexp,omitempty"` // This is a regular rule string used to match and verify that the format is legal.
	Required     bool        `json:"required,omitempty" yaml:"required,omitempty" bson:"required,omitempty"`                // Is the mark mandatory?
}
