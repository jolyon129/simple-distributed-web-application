[![Build Status](https://travis-ci.com/Distributed-Systems-CSGY9223/zl2501-final-project.svg?token=LyWHGctXVCcEk9v6z4HG&branch=master)](https://travis-ci.com/Distributed-Systems-CSGY9223/zl2501-final-project)

# Stage 3 Explanation

I choose `etcd` from `CoreOS`  as the raft implementation. 

Each raft node will expose a RESTful API for DB storage. The `backend` and `auth` service have a new `raftclient` DB engine  which is the abstract DB implementation of the storage for `userstore`,`tweetstore` and `sessionProvider`. Under the hood, the `raftclient` will send requests to raft cluster to get the data. `raftclient` along with `memory` are two different DB engine and can be switched easily by a engine register mechanism(`backend` does not fully support engine register mechanism). 

The `raftclient` is undere the folder `/model/sorage/raftclient` in `backend` and `auth` service. 

In this application, we can have up to 3 cluster nodes. They will expose follwoing API address: 
```go 
    ADDR1 = "http://127.0.0.1:9004"
    ADDR2 = "http://127.0.0.1:9005"
    ADDR3 = "http://127.0.0.1:9006"
```

I document the detailed DB API in the `Postman`, you can check it out through this url: https://documenter.getpostman.com/view/1347930/SzmcbKNh?version=latest

Ideally, all the above addresses will give the same result of DB as they are consistent through raft protocols.

The `raftclient` DB engine, or the wrapper, will try theese API addresses one by one with timeout cancelation till one of them success. And the recently successed one will always be ranked first to try. This gives `backend` and `auth` the `Fault Tolerance` we want.

I also implement a customized trie-tree-based `mux` for the httpserver in `raftcluster`. So that I can easily register RESTful routers. Really having a good time coding on this. 


## Commands

Separately call `make run-raft` `make run-auth`, `make run-backend` and `make run-web` in four terminal sessions
. Then go to `localhost:9000` to enter into the application. (Note: `make run-raft` will only start one raft node and a API server at `localost:9004`. I will demostrate the useage of 3 raft nodes in the final presentation.)

Make Targets:
* `make run-auth`: Start a raft node and expose http server of DB engine at `localhost:9004`
* `make run-auth`: Start the auth server at `localhost:9002`
* `make run-backend`: Start the backend server at `localhost:9001` 
* `make run-web`: Start the web server at `localhost:9000`
* `make test`: Run ginkgo test
* `make build`: Build into `/build` directory
* `make proctoc`: Generate gRPC stubs and distribute into each service directories


## Logic

The logic is same as the stage one. Still, there is some default data.

You can login by using the following usernames or register a new  one.
* User: jolyon129, Password: 123
* User: zl2501, Password: 123


# Stage 2 Explanation

I split the application into three services: web, backend and auth.

The Backend service exposes access to process DB.The auth service authenticates web requests and servers session data in memory. The web service act as a aggregator to composite response from other services.

The web server writes a sessionId in cookie and it communicates to auth server to store the `userId` and `userName` of the current session into the session data layer in the auth server. All processes related to DB are moved into backend server. Hence, the web server is totally stateless. 


## Commands 

Separately call `make run-auth`, `make run-backend` and `make run-web` in three terminal sessions
. Then go to `localhost:9000` to enter into the application.

Please note that you need to start the auth and backend server first, and then start the web server.

Make Targets:
* `make run-auth`: Start the auth server at `localhost:9002`
* `make run-backend`: Start the backend server at `localhost:9001` 
* `make run-web`: Start the web server at `localhost:9000`
* `make test`: Run ginkgo test
* `make build`: Build into `/build` directory.
* `make proctoc`: Generate gRPC stubs and distribute into each service directories.

## Logic

The logic is same as the stage one. Still, there is some default data.

You can login by using the following usernames or register a new  one.
* User: jolyon129, Password: 123
* User: zl2501, Password: 123

## Structure

```
.
├── Makefile    --- makefile
├── README.md
├── autu            --- Auth Service
│   ├── auth.go     --- Entry point of auth service 
│   ├── constant
│   │   └── constant.go
│   ├── go.mod
│   ├── go.sum
│   ├── pb
│   │   └── auth.pb.go	--- Auto generated RPC stub
│   ├── rpc.go          --- Implement RPC calls
│   ├── sessmanager     --- Session manager(call functions in provider)
│   │   ├── const.go
│   │   ├── manager.go
│   │   ├── sessmanager_test.go
│   │   └── sessmanger_suite_test.go
│   └── storage         --- Implement session storage in memory
│       ├── memory      --- Thread-safe memory Implementation
│       │   └── memory.go
│       ├── provider_interface.go   --- Interface for session provider
│       └── session_interface.go    --- Interface for session 
├── backend                 --- Backend Service
│   ├── backend.go          --- Entry point of backend service
│   ├── constant
│   │   └── constant.go
│   ├── go.mod
│   ├── go.sum
│   ├── model               --- Model layer
│   │   ├── model.go
│   │   ├── model_suite_test.go
│   │   ├── model_test.go
│   │   ├── repository      --- Implement Repository
│   │   │   ├── repository_suite_test.go
│   │   │   ├── repository_test.go
│   │   │   ├── tweetrepo.go
│   │   │   └── userrepo.go
│   │   └── storage         --- Implement Storage Layer
│   │       ├── memory      --- Thread-safe memory Implementation
│   │       │   ├── memory.go								
│   │       │   ├── tweetstorage.go
│   │       │   └── userstorage.go
│   │       └── storage_interface.go
│   ├── pb                  --- Auto generated RPC stub
│   │   └── backend.pb.go
│   └── rpc.go              --- Implement RPC calls
├── build                   --- The output of Build
│   ├── auth
│   ├── backend
│   └── web
├── cmd                    --- Commands to execute services(call `StartService` in each service) 
│   ├── auth
│   │   └── auth.go
│   ├── backend
│   │   └── backend.go      --- This also calls `addDefaultData`
│   └── web
│       └── web.go						
├── commonpb                --- Store the proto files
│   ├── auth.proto
│   └── backend.proto
├── go.sum
└── web                     --- Web Service
    ├── constant
    │   └── constant.go
    ├── controller          --- Controller layer
    │   ├── controller.go   --- Other controllers (login,logout, signin, tweet, etc)
    │   ├── home.go         --- Controllers for home page
    │   └── util.go
    ├── go.mod
    ├── go.sum
    ├── middleware          --- Middleware for requests(CheckAuth,Logger,SetHeader)
    │   └── middleware.go
    ├── pb                  --- Auto generated RPC stub and implement Dial
    │   ├── auth.pb.go
    │   ├── backend.pb.go
    │   └── dial.go         --- Dial to auth and backend service
    ├── template            --- HTML files
    │   ├── home.html
    │   ├── index.html
    │   ├── login.html
    │   ├── signup.html
    │   ├── tweet.html
    │   ├── user.html
    │   └── users.html
    └── web.go              --- Entry point and router: route reqeusts to corresponding controllers  

25 directories, 59 files
```


# Stage 1 Explanation

I implemented a memory session layer which is actually LRU under the hood. The memory DB for users and lists is implemented by `list.List` and `set`. 

Every requests need to go throug some middlewares: such as `CheckAuth`, `SetHeader`, `Logger`. The middleware `CheckAUth` act as the `auth` module.

The authtication of requests is done by the middleware `CheckAuth` and the middleware will check the sessionId within the cokie with the session mananger.

## Run the server

WARN: The project is built on Go 1.14 and uses `Go Module`. If your go version is lower than this, the `go mod vendor` may cause errors. 

Please put the project under the `$GOPATH/src` and cd into the project first.

I use `Makefile` to organize commands.
* `make run-web`:  Start the server(will call `go mod vendor` first to download 3rd party packages(Ginkgo and Crypto). The server starts at `http://localhost:9000`
* `make test`: Run `go test -v --race`to call `ginkgo`.
* `make build`: Build `web` into `./build` directory. After building, you can execute `./build/web` to run the server. The working directory has to be the root of the project so that the server can access to the HTML file(which is in `./web/template/`)

## Logic

`make run-web` to start the server.

After starting service, go to `localhost:9000` to enter into the application. You can creat your own user or login by using the predefined test user.

**The uesrname is "zl2501", and the password is "123".**

The user `zl2501` has some tweets and is following the user `jolyon129`(password is 123). After login, the page is redirected to the `home` and you can see your feed(including your own tweets and `jolyon129`'s tweets). I use merge sort to display all tweets chronologically.

Tweet as you want. 

Your can also view other users by clicking `View all users`. On this user list page, you can follow and unfollow others (your feed will change as well).    

## URL 

* `/index`  login or sign up 
* `/home`   view the feed which consists of tweets from the following users (need to check auth ahead. If not login, redirect to the `/index`), and can take other basic actions (logout, tweet, view user list)
* `/users`  display all other users so that the user can follow or unfollow (need to check auth ahead)
* `/user/:username` view some user's tweets


## Project Structure

```
.
├── Makefile        --- makefile
├── README.md       
├── build           --- store the result of building 
├── cmd              
│   └── web        
│       └── web.go  --- call `StartService` and `addDefaultData`
└── web
    ├── auth            
    │   └──  auth.go        --- Authentication Middleware
    │   
    ├── constant            --- Some Configuration and Constants(Port Number, etc)
    │   └── constant.go
    ├── controller          --- Implement controllers for the requests
    │   ├── controller.go   --- Other controllers(`login`,`sigin`,`tweet`, etc)
    │   ├── home.go         --- Seperate file for home controller
    │   └── util.go         --- Utility 
    ├── go.mod              --- Go Module 
    ├── go.sum
    ├── logger              --- Request Logger Middileware
    │   └── logger.go      
    ├── model               --- Implement Model Layer
    │   ├── model.go
    │   ├── model_suite_test.go --- Ginkgo Bootstrap File
    │   ├── model_test.go   -- Tests
    │   ├── repository      --- Implement repository
    │   │   ├── postrepo.go     --- Post/Tweets Repository
    │   │   ├── repository_suite_test.go --- Ginkgo Bootstrap File
    │   │   ├── repository_test.go  --- Tests for userrepo and postrepo
    │   │   └── userrepo.go     --- User Repository
    │   └── storage         --- Implement Storage Layer
    │       ├── memory      --- thread-safe memory Implementation
    │       │   ├── memory.go   --- Register memory implemnetation as the provider of stroage
    │       │   ├── poststorage.go  --- Post/Tweet Storage
    │       │   └── userstorage.go  --- User Storage
    │       └── storage_interface.go    --- Storage interface for users and posts/tweets
    ├── session                 --- Implement session control
    │   ├── provider_interface.go   --- Session provider interface
    │   ├── session_suite_test.go   --- Ginkgo bootstrap file
    │   ├── session_test.go         --- Tests
    │   ├── sessmanager             --- Export session manager to be called by others 
    │   │   ├── const.go        --- session manager configuration
    │   │   └── manager.go      --- Implement Session Manager(Singleton)
    │   └── storage             --- Session Storage. Its LRU under the hood.
    │       ├── memory          
    │       │   └── memory.go   --- thread-safe memory implementation
    │       └── session_interface.go    --- Interaface for session 
    ├── template                --- HTML files
    │   ├── home.html
    │   ├── index.html
    │   ├── login.html
    │   ├── signup.html
    │   ├── tweet.html
    │   ├── user.html
    │   └── users.html
    ├── vendor                 --- The directory for 3rd Party Library(crypto,ginkgo)
    └── web.go                 --- Set up router: route reqeusts to corresponding controllers  
```

# Distributed Systems: Final Project (Requirement)

The project for this course will be to develop (in stages), a distributed, reliable backend in support of a mildly complex (think twitter) social media application.

The development of this project is divided into 3 stages!

## Summary

The stages are split up as follows:

1. A monolithic web application, with all logic and persistence happening on process.  This process exposes a simple sochttptestial media website.
1. The monolithic web application is now split into several services, with a *stateless* web server and *at least one* backend service, which persists data and communicates with the stateless web server via gRPC.
1. The backend service(s) are now stateless, persisting their state in a raft, replicated data store.

These segments are due March 25th, April 8th, and April 29th, respectively.

\newpage

## Stage 1

### Deliverable

#### Summary

A simple web application, comprised of a web server written in Go, serving html files from a single machine.  For this stage of the project, you do not need to persist data in files or a database – instead, keep everything in memory.  Do not use a database for this application.

#### Features

This application needs to have a small number of very clear functions.  The minimum set are as follows:

1. Creating an account, with username and password
1. Logging in as a given user, given username and password
1. Users can follow other users. This action must be reversible.
1. Users can create posts that are associated with their identity.
1. Users can view some feed composed only of content generated by users they follow.

Note that these operations *should* be accomplished efficiently, but *need* not be.  You won't lose points for O(n) algorithms, where an O(logn) solution is possible.  However, you're selling yourself short ;)

If you want to build a system other than this twitter example, that is totally OK.  Just speak to the professor/TA first, to get approval.  We just want to make sure that your application is of comparable difficulty and testability.

#### Tests

You *will* be expected to write tests for your functions!! If you encounter any issues with testing, don't hesitate to ping your TA!

If you want to use fancy testing frameworks (like [ginkgo](https://github.com/onsi/ginkgo)), that's awesome!  Not required, but might be fun.

To be clear we're not setting an arbitrary code coverage metric, but tests are super important, so they'll factor into your team's grade.  If you're nervous about whether your tests are good enough, speak to your TA about it.  We just want to drive home that tests are super important, especially in this phase, as these tests will be your canaries when you split off the backend in step 2.

\newpage

#### Frameworks and Vendoring

On that note, if you want to use frameworks such as buffalo or ginkgo, go ahead, but make sure you practice [vendoring them as dependencies](https://github.com/golang/go/wiki/modules)!  It'll make it easier on your teammates, and easier for me to grade!  That being said, you definitely will not *need* to use any of these frameworks.  As our system is relatively small, they won't save you a ton of time.  It's also worth learning how stl packages like `net/http` and `testing` work, as that's what all of these frameworks are built on.

#### Structure

As friendly advice, you should try to structure your application as follows.  You don't *need* to, but it will definitely make stage 2 easier! This is because you will be decoupling your modules in the second phase, and if you don't structure your project well, like for example, dump the whole code in one or two files, it'll be hard to decouple the modules later. 

        cmd
        |-- web
        |   `-- web.go --> build target (the code that actually runs the server)
        web
        |-- web.go
        |-- config.go --> config for web, determines port to bind to, etc.
        |-- cmd       --> implementation of build target
        `-- auth
            |-- auth.go   --> authentication module, creates a new auth object,
            |                 starts it, stops it.
            |-- config.go --> token validity duration, hash difficulties, etc.
            |-- errors.go --> errors that auth object can reply with
            |-- *.go      --> implement auth; use contexts, generated
            |                 proto, and storage packages
            |-- *_test.go --> don't forget tests!
            |-- storage
            |   |-- storage.go --> storage interface for auth
            |   `-- memory
            |       |-- config.go
            |       `-- memory.go --> *threadsafe* implementation,
                                      using maps and lists
            `-- authpb
                `-- models.proto --> used to generate all data
                                     primitives used by auth module

The auth package is enumerated as an example.  All of auth's sister packages should look pretty similar, and `web` should compose these packages, hiding access to them behind http handlers, as was shown in class.

If you want, you can totally abandon this structure and follow your own path!! It's totally up to you, but please pick something sensible!! :) 

\newpage

### Resources

The following resources might be helpful:

1. This [open source book](https://astaxie.gitbooks.io/build-web-application-with-golang/content/en/), which explains how to build a go web server from fundamentals.  Note that there is way more information in this link than you'll need, but it's a great overview.  Some of the snippets in Chapter 3 might be especially useful.
1. This [blog post](https://www.ardanlabs.com/blog/2017/02/package-oriented-design.html) explains a popular method of organizing projects, which is quite similar to the one your TA demonstrated in-class!
1. ~~This [exploration of Kubernetes source](https://developer.ibm.com/opentech/2017/06/21/tour-kubernetes-source-code-part-one-kubectl-api-server/), which your TA found on the day he was writing this doc, is really cool, and gives you an intimate sense of how the project works.~~
1. Since the above link doesn't work anymore, if you want to learn about Kubernetes, use [this link](https://jvns.ca/blog/2017/06/04/learning-about-kubernetes/).
1. This [example travis script](https://gist.github.com/y0ssar1an/df2dab474520c4086926f672c52db139) for those of you who want to play around with automated builds and tests using Travis!
1. For leveraging contexts, check out [this blog post](https://blog.golang.org/context), which goes over the subject way better than I could!!
1. This might be a bit late to slide in, but [this talk](https://blog.golang.org/advanced-go-concurrency-patterns) on concurrency patterns is great.a
1. For faking http responses, [this package](https://github.com/jarcoal/httpmock) is great.

\newpage

## Stage 2

### Deliverable

In the next stage, you will split off your backend into a separate service.  It can be a monolith (i.e. one service that performs all of the tasks related to your system's state management), or a set of services that each perform one small task, as described in-class.

This service (set) must communicate with the web server using an RPC Framework like gRPC (strongly recommended!) or Thrift.  If you run into any issues getting these frameworks set up, don't hesitate to contact the TA!

At this point, your webpage service (i.e. the service that is rendering templates / serving web pages) should be *stateless*.  In practical terms, you should be able to horizontally scale your web service without causing any concurrency problems on the backend, or UX problems on the front end.  The Web Service should no longer be persisting any data (even session data!!), but instead be fully reliant on the backend service.  If this isn't the case, you still have work to do :)

Depending on how you solved stage 1, this could be REALLY challenging -- you might need to totally re-write your service, in some cases.  If you followed the model your TA showed in-class, however, this might be the least time-intensive stage.  If you have extra time and want a challenge, you might want to implement more features, or split the backend into more distinct services.

### Resources

The following resources might be helpful:

1. [This document](https://grpc.io/docs/quickstart/go.html) is a super concise explanation of installing gRPC and protobufs (the transport serialization used by gRPC).  It also references an example, which is a great resource in and of itself for learning protobuf syntax.  In particular, it contains snippets for creating, registering, and dialling with gRPC servers/clients, which you can directly use in the project.
1. Review the example rpc slides, and the accompanying repo, which are posted on piazza
1. Adam's [gen proto script](https://raw.githubusercontent.com/adamsanghera/go-websub/master/gen-proto-go.sh) :)
1. Your TA! If you have any issues!

## Stage 3

### Summary

In stage 3, we are finally going to give our system a much-needed upgrade -- persistence!!  And not just any kind of persistence -- we're talking about highly-available, replicated persistence!

To accomplish this, we're going to use one of the two popular open source Raft implementations out there, either the one from [CoreOS](https://github.com/etcd-io/etcd/tree/master/raft) or [Hashicorp](https://github.com/hashicorp/raft).  You can choose to use whichever one you prefer; in your final presentation, you to give some rationale for why you picked the one you did! :)

### Deliverable

This is a relatively small step, but it requires learning how your raft implmentation of choice looks like!! It should end up being more of an ops problem than a dev problem.  Again, if you've followed the example architecture, it shouldn't be a huge headache.  If you find that it is, contact your TA!!

You'll need to:

1. Write code that stands up raft nodes (on your local machine is ok!!)
1. Replace your tried and tested storage implementations with a new one, which wraps a raft client for your raft implementation of choice.  This will live alongside your old implementation in the `storage` directory, if you're following the recommended structure.
   1. This involves code that establishes a connection with your raft cluster.  Now, you'll find that your storage implementation's `config`, `start` and `shutdown` actually do some heavy lifting!!

### Details

Each service should have its *own client* wrapper for contacting the raft cluster, but it is totally OK if they are all targeting the same cluster.  If you want, you can target multiple clusters -- good luck getting all those to run on your laptop, though :).

Don't worry about hosting this solution anywhere other than your laptop.  If you want to host it on AWS, heroku, or what have you (just to show off!), that's totally OK.  Just make sure that you have a configuration that will also work on your local machine.  One hack for accomplishing this in the recommended structure, is to build out `cmd` directory as follows:

\newpage

        cmd
        |-- local
        |   |-- web
        |   `-- backend
        `-- aws
            |-- web
            `-- backend

In the above example, your config initializations under `local` would look very different from your config inits under `prod`.  Alternatively, you can read config objects values from environment variables, have different env-setting shell scripts for aws/local, and keep a single `cmd` folder.  It's really up to you!

