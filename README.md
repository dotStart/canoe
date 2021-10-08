Canoe
=====

A small boat which carries your JVM applications across the data stream. Supports Linux, Mac OS and
Windows!

Canoe wraps your Java* applications in a small native wrapper which takes care of locating the right
runtime and tells you if they don't have the right version.

*Supports all JVM compatible languages which may be packaged as self-contained JAR archives.

Usage
-----

All you need to wrap an application for all supported operating systems is one simple command:

```
canoegen wrap -in my.jar -out bin -runtime-version 16
```

For more customization options, refer to `canoegen help wrap`!

License
-------

```
Copyright [yyyy] [name of copyright owner]

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```
