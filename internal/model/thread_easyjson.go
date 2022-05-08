// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package model

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson2d00218DecodeGithubComPerlinleoTechnoparkMailRuForumDatabaseInternalModel(in *jlexer.Lexer, out *ThreadUpdate) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "message":
			out.Message = string(in.String())
		case "title":
			out.Title = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson2d00218EncodeGithubComPerlinleoTechnoparkMailRuForumDatabaseInternalModel(out *jwriter.Writer, in ThreadUpdate) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"message\":"
		out.RawString(prefix[1:])
		out.String(string(in.Message))
	}
	{
		const prefix string = ",\"title\":"
		out.RawString(prefix)
		out.String(string(in.Title))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ThreadUpdate) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson2d00218EncodeGithubComPerlinleoTechnoparkMailRuForumDatabaseInternalModel(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ThreadUpdate) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson2d00218EncodeGithubComPerlinleoTechnoparkMailRuForumDatabaseInternalModel(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ThreadUpdate) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson2d00218DecodeGithubComPerlinleoTechnoparkMailRuForumDatabaseInternalModel(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ThreadUpdate) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson2d00218DecodeGithubComPerlinleoTechnoparkMailRuForumDatabaseInternalModel(l, v)
}
func easyjson2d00218DecodeGithubComPerlinleoTechnoparkMailRuForumDatabaseInternalModel1(in *jlexer.Lexer, out *Thread) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "author":
			out.Author = string(in.String())
		case "created":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.Created).UnmarshalJSON(data))
			}
		case "forum":
			out.Forum = string(in.String())
		case "id":
			out.ID = int32(in.Int32())
		case "message":
			out.Message = string(in.String())
		case "slug":
			out.Slug = string(in.String())
		case "title":
			out.Title = string(in.String())
		case "votes":
			out.Votes = int32(in.Int32())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson2d00218EncodeGithubComPerlinleoTechnoparkMailRuForumDatabaseInternalModel1(out *jwriter.Writer, in Thread) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"author\":"
		out.RawString(prefix[1:])
		out.String(string(in.Author))
	}
	{
		const prefix string = ",\"created\":"
		out.RawString(prefix)
		out.Raw((in.Created).MarshalJSON())
	}
	{
		const prefix string = ",\"forum\":"
		out.RawString(prefix)
		out.String(string(in.Forum))
	}
	if in.ID != 0 {
		const prefix string = ",\"id\":"
		out.RawString(prefix)
		out.Int32(int32(in.ID))
	}
	{
		const prefix string = ",\"message\":"
		out.RawString(prefix)
		out.String(string(in.Message))
	}
	{
		const prefix string = ",\"slug\":"
		out.RawString(prefix)
		out.String(string(in.Slug))
	}
	{
		const prefix string = ",\"title\":"
		out.RawString(prefix)
		out.String(string(in.Title))
	}
	{
		const prefix string = ",\"votes\":"
		out.RawString(prefix)
		out.Int32(int32(in.Votes))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Thread) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson2d00218EncodeGithubComPerlinleoTechnoparkMailRuForumDatabaseInternalModel1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Thread) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson2d00218EncodeGithubComPerlinleoTechnoparkMailRuForumDatabaseInternalModel1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Thread) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson2d00218DecodeGithubComPerlinleoTechnoparkMailRuForumDatabaseInternalModel1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Thread) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson2d00218DecodeGithubComPerlinleoTechnoparkMailRuForumDatabaseInternalModel1(l, v)
}
func easyjson2d00218DecodeGithubComPerlinleoTechnoparkMailRuForumDatabaseInternalModel2(in *jlexer.Lexer, out *NewThread) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "author":
			out.Author = string(in.String())
		case "created":
			out.Created = string(in.String())
		case "forum":
			out.Forum = string(in.String())
		case "message":
			out.Message = string(in.String())
		case "slug":
			out.Slug = string(in.String())
		case "title":
			out.Title = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson2d00218EncodeGithubComPerlinleoTechnoparkMailRuForumDatabaseInternalModel2(out *jwriter.Writer, in NewThread) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"author\":"
		out.RawString(prefix[1:])
		out.String(string(in.Author))
	}
	{
		const prefix string = ",\"created\":"
		out.RawString(prefix)
		out.String(string(in.Created))
	}
	{
		const prefix string = ",\"forum\":"
		out.RawString(prefix)
		out.String(string(in.Forum))
	}
	{
		const prefix string = ",\"message\":"
		out.RawString(prefix)
		out.String(string(in.Message))
	}
	{
		const prefix string = ",\"slug\":"
		out.RawString(prefix)
		out.String(string(in.Slug))
	}
	{
		const prefix string = ",\"title\":"
		out.RawString(prefix)
		out.String(string(in.Title))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v NewThread) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson2d00218EncodeGithubComPerlinleoTechnoparkMailRuForumDatabaseInternalModel2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v NewThread) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson2d00218EncodeGithubComPerlinleoTechnoparkMailRuForumDatabaseInternalModel2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *NewThread) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson2d00218DecodeGithubComPerlinleoTechnoparkMailRuForumDatabaseInternalModel2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *NewThread) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson2d00218DecodeGithubComPerlinleoTechnoparkMailRuForumDatabaseInternalModel2(l, v)
}