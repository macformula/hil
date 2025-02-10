# 1. Creating the `main` application

## Introduction
The `macfe/hil` package is a software-in-the-loop testing package developed in Go.
This tutorial will walkthrough the process of using this package to set up SIL tests. 

## Starting our application
To start, create an empty directory to hold all of your SIL related files. We will call it `sil`. Set up the `sil` directory to have this file structure. 

```
sil
├─ config
├─ cmd
├─ ecu
├─ pinout
├─ state
└─ results
```