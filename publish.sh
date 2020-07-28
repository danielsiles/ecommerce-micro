#!/bin/bash

git add . && git commit -m "updating"
git tag -a v0.0.3 -m "updating tag"
git push origin master --tags
