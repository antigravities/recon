# recon
Report system metrics to AWS CloudWatch

## Install

### Configuration
Configuration is done via environment variables.

You can specify a `.env` file in the same working directory, or manually specify your environment variables via other methods. Make sure you have an AWS config defined, either via `~/.aws/config`, via environment variables directly, or via `.env`.

```
RE_LOG_GROUP=your_log_group
RE_NAMESPACE=your_cw_namespace
```

### Run once, i.e. via `cron`
```
GO111MODULE=on go get -u get.cutie.cafe/recon # or download a release
# add to your crontab
$GOPATH/bin/recon
```

### Run as a daemon
```
GO111MODULE=on go get -u get.cutie.cafe/recon # or download a release
$GOPATH/bin/recon -daemon
```

## License
```
recon
Copyright (C) 2021 Alexandra Frock

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.