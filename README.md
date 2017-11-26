# fallingstar
[![Go Report Card](https://goreportcard.com/badge/github.com/clly/fallingstar)](https://goreportcard.com/report/github.com/clly/fallingstar)
[![Build Status](https://travis-ci.org/clly/fallingstar.svg?branch=travis)](https://travis-ci.org/clly/fallingstar)
[![Godoc](https://godoc.org/github.com/shoenig/toolkit?status.svg)](https://godoc.org/github.com/clly/fallingstar)
[![License](https://img.shields.io/github/license/clly/fallingstar.svg)](LICENSE)

Automatically download starred repositories into owner/repo namespace
folders. Fallingstar was created as protection against a potential
["left-pad" issue](https://arstechnica.com/information-technology/2016/03/rage-quit-coder-unpublished-17-lines-of-javascript-and-broke-the-internet/).
You should always ensure that you have a LICENSE (open source or otherwise)
for whatever software you're pulling out of github. I'm not a lawyer and
fallingstar does not attempt to validate LICENSE text.

## Short version

NPM gave some company the name of a package that another developer had already
registered. The dev got mad and deleted all of his packages from NPM. One of
those packages was left-pad which was used to right justify text. It was also
used by a bunch of javascript build tools and a bunch of stuff started breaking.
Someone eventually published a functional equivalent of left-pad but it still
broke a bunch of stuff

## Why does this exist

This is an attempt to ensure that code that at some point lived in the
"public domain" ([techincally all rights reserved](https://www.infoworld.com/article/2615869/open-source-software/github-needs-to-take-open-source-seriously.html]))
will always exist for you locally. It's not as good as republishing it
elsewhere but it's better than nothing.

## Can I use this for project "X"

Again I'm not a lawyer. You need to ensure that you can use the project
and if it has a LICENSE
