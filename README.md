packer-post-processor-virtualbox
================================

Deploy packer artifacts to remote virtualbox host

Borrows heavily from mheidenr's work on a VMWare post-processor @ https://github.com/mheidenr/packer-post-processor-ovftool

Dependancies
------------
 ssh
 scp
 
Build
-----
    go plugin\post-processor-virtualbox\main.go
