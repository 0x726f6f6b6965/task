#!/usr/bin bash
echo TASK_VERSION task-svc:$(git describe --tags --abbrev=0)
echo task-svc:$(git describe --tags --abbrev=0) > api/tag.txt