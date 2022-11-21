(window.webpackJsonp=window.webpackJsonp||[]).push([[170],{807:function(e,t,n){"use strict";n.r(t);var i=n(1),a=Object(i.a)({},(function(){var e=this,t=e.$createElement,n=e._self._c||t;return n("ContentSlotsDistributor",{attrs:{"slot-key":e.$parent.slotKey}},[n("h1",{attrs:{id:"rfc-023-semi-permanent-testnet"}},[n("a",{staticClass:"header-anchor",attrs:{href:"#rfc-023-semi-permanent-testnet"}},[e._v("#")]),e._v(" RFC 023: Semi-permanent Testnet")]),e._v(" "),n("h2",{attrs:{id:"changelog"}},[n("a",{staticClass:"header-anchor",attrs:{href:"#changelog"}},[e._v("#")]),e._v(" Changelog")]),e._v(" "),n("ul",[n("li",[e._v("2022-07-28: Initial draft (@mark-rushakoff)")]),e._v(" "),n("li",[e._v("2022-07-29: Renumber to 023, minor clarifications (@mark-rushakoff)")])]),e._v(" "),n("h2",{attrs:{id:"abstract"}},[n("a",{staticClass:"header-anchor",attrs:{href:"#abstract"}},[e._v("#")]),e._v(" Abstract")]),e._v(" "),n("p",[e._v("This RFC discusses a long-lived testnet, owned and operated by the Tendermint engineers.\nBy owning and operating a production-like testnet,\nthe team who develops Tendermint becomes more capable of discovering bugs that\nonly arise in production-like environments.\nThey also build expertise in operating Tendermint;\nthis will help guide the development of Tendermint towards operator-friendly design.")]),e._v(" "),n("p",[e._v("The RFC details a rough roadmap towards a semi-permanent testnet, some of the considered tradeoffs,\nand the expected outcomes from following this roadmap.")]),e._v(" "),n("h2",{attrs:{id:"background"}},[n("a",{staticClass:"header-anchor",attrs:{href:"#background"}},[e._v("#")]),e._v(" Background")]),e._v(" "),n("p",[e._v("The author's understanding -- which is limited as a new contributor to the Tendermint project --\nis that Tendermint development has been largely treated as a library for other projects to consume.\nOf course effort has been spent on unit tests, end-to-end tests, and integration tests.\nBut whether developing a library or an application,\nthere is no substitute for putting the software under a production-like load.")]),e._v(" "),n("p",[e._v('First, there are classes of bugs that are unrealistic to discover in environments\nthat do not resemble production.\nBut perhaps more importantly, there are "operational features" that are best designed\nby the authors of a given piece of software.\nFor instance, does the software have sufficient observability built-in?\nAre the reported metrics useful?\nAre the log messages clear and sufficiently detailed, without being too noisy?')]),e._v(" "),n("p",[e._v("Furthermore, if the library authors are not only building --\nbut also maintaining and operating -- an application built on top of their library,\nthe authors will have a greatly increased confidence that their library's API\nis appropriate for other application authors.")]),e._v(" "),n("p",[e._v("Once the decision has been made to run and operate a service,\none of the next strategic questions is that of deploying said service.\nThe author strongly holds the opinion that, when possible,\na continuous delivery model offers the most compelling set of advantages:")]),e._v(" "),n("ul",[n("li",[e._v("The code on a particular branch (likely "),n("code",[e._v("main")]),e._v(" or "),n("code",[e._v("master")]),e._v(") is exactly what is,\nor what will very soon be, running in production")]),e._v(" "),n("li",[e._v("There are no manual steps involved in deploying -- other than merging your pull request,\nwhich you had to do anyway")]),e._v(" "),n("li",[e._v("A bug discovered in production can be rapidly confirmed as fixed in production")])]),e._v(" "),n("p",[e._v("In summary, if the tendermint authors build, maintain, and continuously deliver an application\nintended to serve as a long-lived testnet, they will be able to state with confidence:")]),e._v(" "),n("ul",[n("li",[e._v("We operate the software in a production-like environment and we have observed it to be\nstable and performant to our requirements")]),e._v(" "),n("li",[e._v("We have discovered issues in production before any external parties have consumed our software,\nand we have addressed said issues")]),e._v(" "),n("li",[e._v("We have successfully used the observability tooling built into our software\n(perhaps in conjunction with other off-the-shelf tooling)\nto diagnose and debug issues in production")])]),e._v(" "),n("h2",{attrs:{id:"discussion"}},[n("a",{staticClass:"header-anchor",attrs:{href:"#discussion"}},[e._v("#")]),e._v(" Discussion")]),e._v(" "),n("p",[e._v("The Discussion Section proposes a variety of aspects of maintaining a testnet for Tendermint.")]),e._v(" "),n("h3",{attrs:{id:"number-of-testnets"}},[n("a",{staticClass:"header-anchor",attrs:{href:"#number-of-testnets"}},[e._v("#")]),e._v(" Number of testnets")]),e._v(" "),n("p",[e._v("There should probably be one testnet per maintained branch of Tendermint,\ni.e. one for the "),n("code",[e._v("main")]),e._v(" branch\nand one per "),n("code",[e._v("v0.N.x")]),e._v(" branch that the authors maintain.")]),e._v(" "),n("p",[e._v("There may also exist testnets for long-lived feature branches.")]),e._v(" "),n("p",[e._v("We may eventually discover that there is good reason to run more than one testnet for a branch,\nperhaps due to a significant configuration variation.")]),e._v(" "),n("h3",{attrs:{id:"testnet-lifecycle"}},[n("a",{staticClass:"header-anchor",attrs:{href:"#testnet-lifecycle"}},[e._v("#")]),e._v(" Testnet lifecycle")]),e._v(" "),n("p",[e._v('The document has used the terms "long-lived" and "semi-permanent" somewhat interchangeably.\nThe intent of the testnet being discussed in this RFC is to exist indefinitely;\nbut there is a practical understanding that there will be testnet instances\nwhich will be retired due to a variety of reasons.\nFor instance, once a release branch is no longer supported,\nits corresponding testnet should be torn down.')]),e._v(" "),n("p",[e._v("In general, new commits to branches with corresponding testnets\nshould result in an in-place upgrade of all nodes in the testnet\nwithout any data loss and without requiring new configuration.\nThe mechanism for achieving this is outside the scope of this RFC.")]),e._v(" "),n("p",[e._v("However, it is also expected that there will be\nbreaking changes during the development of the "),n("code",[e._v("main")]),e._v(" branch.\nFor instance, suppose there is an unreleased feature involving storage on disk,\nand the developers need to change the storage format.\nIt should be at the developers' discretion whether it is feasible and worthwhile\nto introduce an intermediate commit that translates the old format to the new format,\nor if it would be preferable to just destroy the testnet and start from scratch\nwithout any data in the old format.")]),e._v(" "),n("p",[e._v("Similarly, if a developer inadvertently pushed a breaking change to an unreleased feature,\nthey are free to make a judgement call between reverting the change,\nadding a commit to allow a forward migration,\nor simply forcing the testnet to recreate.")]),e._v(" "),n("h3",{attrs:{id:"testnet-maintenance-investment"}},[n("a",{staticClass:"header-anchor",attrs:{href:"#testnet-maintenance-investment"}},[e._v("#")]),e._v(" Testnet maintenance investment")]),e._v(" "),n("p",[e._v("While there is certainly engineering effort required to build the tooling and infrastructure\nto get the testnets up and running,\nthe intent is that a running testnet requires no manual upkeep under normal conditions.")]),e._v(" "),n("p",[e._v("It is expected that a subset of the Tendermint engineers are familiar with and engaged in\nwriting the software to maintain and build the testnet infrastructure,\nbut the rest of the team should not need any involvement in authoring that code.")]),e._v(" "),n("p",[e._v("The testnets should be configured to send notifications for events requiring triage,\nsuch as a chain halt or a node OOMing.\nThe time investment necessary to address the underlying issues for those kind of events\nis unpredictable.")]),e._v(" "),n("p",[e._v("Aside from triaging exceptional events, an engineer may choose to spend some time\ncollecting metrics or profiles from testnet nodes to check performance details\nbefore and after a particular change;\nor they may inspect logs associated with an expected behavior change.\nBut during day-to-day work, engineers are not expected to spend any considerable time\ndirectly interacting with the testnets.")]),e._v(" "),n("p",[e._v("If we discover that there are any routine actions engineers must take against the testnet\nthat take any substantial focused time,\nthose actions should be automated to a one-line command as much as is reasonable.")]),e._v(" "),n("h3",{attrs:{id:"testnet-mvp"}},[n("a",{staticClass:"header-anchor",attrs:{href:"#testnet-mvp"}},[e._v("#")]),e._v(" Testnet MVP")]),e._v(" "),n("p",[e._v("The minimum viable testnet meets this set of features:")]),e._v(" "),n("ul",[n("li",[e._v("The testnet self-updates following a new commit pushed to Tendermint's "),n("code",[e._v("main")]),e._v(" branch on GitHub\n(there are some omitted steps here, such as CI building appropriate binaries and\nsomehow notifying the testnet that a new build is available)")]),e._v(" "),n("li",[e._v("The testnet runs the Tendermint KV store for MVP")]),e._v(" "),n("li",[e._v("The testnet operators are notified if:\n"),n("ul",[n("li",[e._v("Any node's process exits for any reason other than a restart for a new binary")]),e._v(" "),n("li",[e._v("Any node stops updating blocks, and by extension if a chain halt occurs")]),e._v(" "),n("li",[e._v("No other observability will be considered for MVP")])])]),e._v(" "),n("li",[e._v("The testnet has a minimum of 1 full node and 3 validators")]),e._v(" "),n("li",[e._v("The testnet has a reasonably low, constant throughput of transactions -- say 30 tx/min --\nand the testnet operators are notified if that throughput drops below 75% of target\nsustained over 5 minutes")]),e._v(" "),n("li",[e._v("The testnet only needs to run in a single datacenter/cloud-region for MVP,\ni.e. running in multiple datacenters is out of scope for MVP")]),e._v(" "),n("li",[e._v("The testnet is running directly on VMs or compute instances;\nwhile Kubernetes or other orchestration frameworks may offer many significant advantages,\nthe Tendermint engineers should not be required to learn those tools in order to\nperform basic debugging")])]),e._v(" "),n("h3",{attrs:{id:"testnet-medium-term-goals"}},[n("a",{staticClass:"header-anchor",attrs:{href:"#testnet-medium-term-goals"}},[e._v("#")]),e._v(" Testnet medium-term goals")]),e._v(" "),n("p",[e._v("The medium-term goals are intended to be achievable within the 6-12 month time range\nfollowing the launch of MVP.\nThese goals could realistically be roadmapped following the launch of the MVP testnet.")]),e._v(" "),n("ul",[n("li",[e._v("The "),n("code",[e._v("main")]),e._v(" testnet has more than 20 nodes (completely arbitrary -- 5x more than 1+3 at MVP)")]),e._v(" "),n("li",[e._v("In addition to the "),n("code",[e._v("main")]),e._v(" testnet,\nthere is at least one testnet associated with one release branch")]),e._v(" "),n("li",[e._v("The testnet no longer is simply running the Tendermint KV store;\nnow it is built on a more complex, custom application\nthat deliberately exercises a greater portion of the Tendermint stack")]),e._v(" "),n("li",[e._v('Each testnet is spread across at least two cloud providers,\nin order to communicate over a network more closely resembling use of Tendermint in "real" chains')]),e._v(" "),n("li",[e._v('The node updates have some "jitter",\nwith some nodes updating immediately when a new build is available,\nand others delaying up to perhaps 30-60 minutes')]),e._v(" "),n("li",[e._v("The team has published some form of dashboards that have served well for debugging,\nwhich external parties can copy/modify to their needs\n"),n("ul",[n("li",[e._v("The dashboards must include metrics published by Tendermint nodes;\nthere should be both OS- or runtime-level metrics such as memory in use,\nand application-level metrics related to the underlying blockchain")]),e._v(" "),n("li",[e._v('"Published" in this context is more in the spirit of "shared with the community",\nnot "produced a supported open source tool" --\nthis could be published to GitHub with a warning that no support is offered,\nor it could simply be a blog post detailing what has worked for the Tendermint developers')]),e._v(" "),n("li",[e._v("The dashboards will likely be implemented on free and open source tooling,\nbut that is not a hard requirement if paid software is more appropriate")])])]),e._v(" "),n("li",[e._v("The team has produced a reference model of a log aggregation stack that external parties can use\n"),n("ul",[n("li",[e._v('Similar to the "published" dashboards, this only needs to be "shared" rather than "supported"')])])]),e._v(" "),n("li",[e._v("Chaos engineering has begun being integrated into the testnets\n(this could be periodic CPU limiting or deliberate network interference, etc.\nbut it probably would not be filesystem corruption)")]),e._v(" "),n("li",[e._v("Each testnet has at least one node running a build with the Go race detector enabled")]),e._v(" "),n("li",[e._v("The testnet contains some kind of generalized notification system built in:\n"),n("ul",[n("li",[e._v('Tendermint code grows "watchdog" systems built in to validate things like\nsubsystems have not deadlocked; e.g. if the watchdog can\'t acquire and immediately release\na particular mutex once in every 5-minute period, it is near certain that the target\nsubsystem has deadlocked, and an alert must be sent to the engineering team.\n(Outside of the testnet, the watchdogs could be disabled, or they could panic on failure.)')]),e._v(" "),n("li",[e._v("The notification system does some deduplication to minimize spam on system failure")])])])]),e._v(" "),n("h3",{attrs:{id:"testnet-long-term-vision"}},[n("a",{staticClass:"header-anchor",attrs:{href:"#testnet-long-term-vision"}},[e._v("#")]),e._v(" Testnet long-term vision")]),e._v(" "),n("p",[e._v("The long-term vision includes goals that are not necessary for short- or medium-term success,\nbut which would support building an increasingly stable and performant product.\nThese goals would generally be beyond the one-year plan,\nand therefore they would not be part of initial planning.")]),e._v(" "),n("ul",[n("li",[e._v("There is a centralized dashboard to get a quick overview of all testnets,\nor at least one centralized dashboard per testnet,\nshowing TBD basic information")]),e._v(" "),n("li",[e._v("Testnets include cloud spot instances which periodically and abruptly join and leave the network")]),e._v(" "),n("li",[e._v("The testnets are a heterogeneous mixture of straight VMs and Docker containers,\nthereby more closely representing production blockchains")]),e._v(" "),n("li",[e._v("Testnets have some manner of continuous profiling,\nso that we can produce an apples-to-apples comparison of CPU/memory cost of particular operations")])]),e._v(" "),n("h3",{attrs:{id:"testnet-non-goals"}},[n("a",{staticClass:"header-anchor",attrs:{href:"#testnet-non-goals"}},[e._v("#")]),e._v(" Testnet non-goals")]),e._v(" "),n("p",[e._v("There are some things we are explicitly not trying to achieve with long-lived testnets:")]),e._v(" "),n("ul",[n("li",[e._v("The Tendermint engineers will NOT be responsible for the testnets' availability\noutside of working hours; there will not be any kind of on-call schedule")]),e._v(" "),n("li",[e._v("As a result of the 8x5 support noted in the previous point,\nthere will be NO guarantee of uptime or availability for any testnet")]),e._v(" "),n("li",[e._v("The testnets will NOT be used to gate pull requests;\nthat responsibility belongs to unit tests, end-to-end tests, and integration tests")]),e._v(" "),n("li",[e._v("Similarly, the testnet will NOT be used to automate any changes back into Tendermint source code;\nwe will not automatically create a revert commit due to a failed rollout, for instance")]),e._v(" "),n("li",[e._v("The testnets are NOT intended to have participation from machines outside of the\nTendermint engineering team's control, as the Tendermint engineers are expected\nto have full access to any instance where they may need to debug an issue")]),e._v(" "),n("li",[e._v('While there will certainly be individuals within the Tendermint engineering team\nwho will continue to build out their individual "devops" skills to produce\nthe infrastructure for the testnet, it is NOT a goal that every Tendermint engineer\nis even '),n("em",[e._v("familiar")]),e._v(" with the tech stack involved, whether it is Ansible, Terraform,\nKubernetes, etc.\nAs a rule of thumb, all engineers should be able to get shell access on any given instance\nand should have access to the instance's logs.\nLittle if any further operational skills will be expected.")]),e._v(" "),n("li",[e._v("The testnets are not intended to be "),n("em",[e._v("created")]),e._v(' for one-off experiments.\nWhile there is nothing wrong with an engineer directly interacting with a testnet\nto try something out,\na testnet comes with a considerable amount of "baggage", so end-to-end or integration tests\nare closer to the intent for "trying something to see what happens".\nDirect interaction should be limited to standard blockchain operations,\n'),n("em",[e._v("not")]),e._v(" modifying configuration of nodes.")]),e._v(" "),n("li",[e._v('Likewise, the purpose of the testnet is not to run specific "tests" per se,\nbut rather to demonstrate that Tendermint blockchains as a whole are stable\nunder a production load.\nOf course we will inject faults periodically, but the intent is to observe and prove that\nthe testnet is resilient to those faults.\nIt would be the responsibility of a lower-level test to demonstrate e.g.\nthat the network continues when a single validator disappears without warning.')]),e._v(" "),n("li",[e._v("The testnet descriptions in this document are scoped only to building directly on Tendermint;\nintegrating with the Cosmos SDK, or any other third-party library, is out of scope")])]),e._v(" "),n("h3",{attrs:{id:"team-outcomes-as-a-result-of-maintaining-and-operating-a-testnet"}},[n("a",{staticClass:"header-anchor",attrs:{href:"#team-outcomes-as-a-result-of-maintaining-and-operating-a-testnet"}},[e._v("#")]),e._v(" Team outcomes as a result of maintaining and operating a testnet")]),e._v(" "),n("p",[e._v("Finally, this section reiterates what team growth we expect by running semi-permanent testnets.")]),e._v(" "),n("ul",[n("li",[e._v("Confidence that Tendermint is stable under a particular production-like load")]),e._v(" "),n("li",[e._v("Familiarity with typical production behavior of Tendermint, e.g. what the logs look like,\nwhat the memory footprint looks like, and what kind of throughput is reasonable\nfor a network of a particular size")]),e._v(" "),n("li",[e._v("Comfort and familiarity in manually inspecting a misbehaving or failing node")]),e._v(" "),n("li",[e._v("Confidence that Tendermint ships sufficient tooling for external users\nto operate their nodes")]),e._v(" "),n("li",[e._v("Confidence that Tendermint exposes useful metrics, and comfort interpreting those metrics")]),e._v(" "),n("li",[e._v("Produce useful reference documentation that gives operators confidence to run Tendermint nodes")])])])}),[],!1,null,null,null);t.default=a.exports}}]);