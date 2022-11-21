(window.webpackJsonp=window.webpackJsonp||[]).push([[138],{744:function(e,t,a){"use strict";a.r(t);var o=a(1),r=Object(o.a)({},(function(){var e=this,t=e.$createElement,a=e._self._c||t;return a("ContentSlotsDistributor",{attrs:{"slot-key":e.$parent.slotKey}},[a("h1",{attrs:{id:"adr-74-migrate-timeout-parameters-to-consensus-parameters"}},[a("a",{staticClass:"header-anchor",attrs:{href:"#adr-74-migrate-timeout-parameters-to-consensus-parameters"}},[e._v("#")]),e._v(" ADR 74: Migrate Timeout Parameters to Consensus Parameters")]),e._v(" "),a("h2",{attrs:{id:"changelog"}},[a("a",{staticClass:"header-anchor",attrs:{href:"#changelog"}},[e._v("#")]),e._v(" Changelog")]),e._v(" "),a("ul",[a("li",[e._v("03-Jan-2022: Initial draft (@williambanfield)")]),e._v(" "),a("li",[e._v("13-Jan-2022: Updated to indicate work on upgrade path needed (@williambanfield)")])]),e._v(" "),a("h2",{attrs:{id:"status"}},[a("a",{staticClass:"header-anchor",attrs:{href:"#status"}},[e._v("#")]),e._v(" Status")]),e._v(" "),a("p",[e._v("Proposed")]),e._v(" "),a("h2",{attrs:{id:"context"}},[a("a",{staticClass:"header-anchor",attrs:{href:"#context"}},[e._v("#")]),e._v(" Context")]),e._v(" "),a("h3",{attrs:{id:"background"}},[a("a",{staticClass:"header-anchor",attrs:{href:"#background"}},[e._v("#")]),e._v(" Background")]),e._v(" "),a("p",[e._v("Tendermint's consensus timeout parameters are currently configured locally by each validator\nin the validator's "),a("a",{attrs:{href:"https://github.com/tendermint/tendermint/blob/5cc980698a3402afce76b26693ab54b8f67f038b/config/toml.go#L425-L440",target:"_blank",rel:"noopener noreferrer"}},[e._v("config.toml"),a("OutboundLink")],1),e._v(".\nThis means that the validators on a Tendermint network may have different timeouts\nfrom each other. There is no reason for validators on the same network to configure\ndifferent timeout values. Proper functioning of the Tendermint consensus algorithm\nrelies on these parameters being uniform across validators.")]),e._v(" "),a("p",[e._v("The configurable values are as follows:")]),e._v(" "),a("ul",[a("li",[a("code",[e._v("TimeoutPropose")]),e._v(" "),a("ul",[a("li",[e._v("How long the consensus algorithm waits for a proposal block before issuing a prevote.")]),e._v(" "),a("li",[e._v("If no prevote arrives by "),a("code",[e._v("TimeoutPropose")]),e._v(", then the consensus algorithm will issue a nil prevote.")])])]),e._v(" "),a("li",[a("code",[e._v("TimeoutProposeDelta")]),e._v(" "),a("ul",[a("li",[e._v("How much the "),a("code",[e._v("TimeoutPropose")]),e._v(" grows each round.")])])]),e._v(" "),a("li",[a("code",[e._v("TimeoutPrevote")]),e._v(" "),a("ul",[a("li",[e._v("How long the consensus algorithm waits after receiving +2/3 prevotes with\nno quorum for a value before issuing a precommit for nil.\n(See the "),a("a",{attrs:{href:"https://arxiv.org/pdf/1807.04938.pdf",target:"_blank",rel:"noopener noreferrer"}},[e._v("arXiv paper"),a("OutboundLink")],1),e._v(", Algorithm 1, Line 34)")])])]),e._v(" "),a("li",[a("code",[e._v("TimeoutPrevoteDelta")]),e._v(" "),a("ul",[a("li",[e._v("How much the "),a("code",[e._v("TimeoutPrevote")]),e._v(" increases with each round.")])])]),e._v(" "),a("li",[a("code",[e._v("TimeoutPrecommit")]),e._v(" "),a("ul",[a("li",[e._v("How long the consensus algorithm waits after receiving +2/3 precommits that\ndo not have a quorum for a value before entering the next round.\n(See the "),a("a",{attrs:{href:"https://arxiv.org/pdf/1807.04938.pdf",target:"_blank",rel:"noopener noreferrer"}},[e._v("arXiv paper"),a("OutboundLink")],1),e._v(", Algorithm 1, Line 47)")])])]),e._v(" "),a("li",[a("code",[e._v("TimeoutPrecommitDelta")]),e._v(" "),a("ul",[a("li",[e._v("How much the "),a("code",[e._v("TimeoutPrecommit")]),e._v(" increases with each round.")])])]),e._v(" "),a("li",[a("code",[e._v("TimeoutCommit")]),e._v(" "),a("ul",[a("li",[e._v("How long the consensus algorithm waits after committing a block but before starting the new height.")]),e._v(" "),a("li",[e._v("This gives a validator a chance to receive slow precommits.")])])]),e._v(" "),a("li",[a("code",[e._v("SkipTimeoutCommit")]),e._v(" "),a("ul",[a("li",[e._v("Make progress as soon as the node has 100% of the precommits.")])])])]),e._v(" "),a("h3",{attrs:{id:"overview-of-change"}},[a("a",{staticClass:"header-anchor",attrs:{href:"#overview-of-change"}},[e._v("#")]),e._v(" Overview of Change")]),e._v(" "),a("p",[e._v("We will consolidate the timeout parameters and migrate them from the node-local\n"),a("code",[e._v("config.toml")]),e._v(" file into the network-global consensus parameters.")]),e._v(" "),a("p",[e._v("The 8 timeout parameters will be consolidated down to 6. These will be as follows:")]),e._v(" "),a("ul",[a("li",[a("code",[e._v("TimeoutPropose")]),e._v(" "),a("ul",[a("li",[e._v("Same as current "),a("code",[e._v("TimeoutPropose")]),e._v(".")])])]),e._v(" "),a("li",[a("code",[e._v("TimeoutProposeDelta")]),e._v(" "),a("ul",[a("li",[e._v("Same as current "),a("code",[e._v("TimeoutProposeDelta")]),e._v(".")])])]),e._v(" "),a("li",[a("code",[e._v("TimeoutVote")]),e._v(" "),a("ul",[a("li",[e._v("How long validators wait for votes in both the prevote\nand precommit phase of the consensus algorithm. This parameter subsumes\nthe current "),a("code",[e._v("TimeoutPrevote")]),e._v(" and "),a("code",[e._v("TimeoutPrecommit")]),e._v(" parameters.")])])]),e._v(" "),a("li",[a("code",[e._v("TimeoutVoteDelta")]),e._v(" "),a("ul",[a("li",[e._v("How much the "),a("code",[e._v("TimeoutVote")]),e._v(" will grow each successive round.\nThis parameter subsumes the current "),a("code",[e._v("TimeoutPrevoteDelta")]),e._v(" and "),a("code",[e._v("TimeoutPrecommitDelta")]),e._v("\nparameters.")])])]),e._v(" "),a("li",[a("code",[e._v("TimeoutCommit")]),e._v(" "),a("ul",[a("li",[e._v("Same as current "),a("code",[e._v("TimeoutCommit")]),e._v(".")])])]),e._v(" "),a("li",[a("code",[e._v("BypassCommitTimeout")]),e._v(" "),a("ul",[a("li",[e._v("Same as current "),a("code",[e._v("SkipTimeoutCommit")]),e._v(", renamed for clarity.")])])])]),e._v(" "),a("p",[e._v("A safe default will be provided by Tendermint for each of these parameters and\nnetworks will be able to update the parameters as they see fit. Local updates\nto these parameters will no longer be possible; instead, the application will control\nupdating the parameters. Applications using the Cosmos SDK will be automatically be\nable to change the values of these consensus parameters "),a("a",{attrs:{href:"https://github.com/cosmos/cosmos-sdk/issues/6197",target:"_blank",rel:"noopener noreferrer"}},[e._v("via a governance proposal"),a("OutboundLink")],1),e._v(".")]),e._v(" "),a("p",[e._v("This change is low-risk. While parameters are locally configurable, many running chains\ndo not change them from their default values. For example, initializing\na node on Osmosis, Terra, and the Cosmos Hub using the their "),a("code",[e._v("init")]),e._v(" command produces\na "),a("code",[e._v("config.toml")]),e._v(" with Tendermint's default values for these parameters.")]),e._v(" "),a("h3",{attrs:{id:"why-this-parameter-consolidation"}},[a("a",{staticClass:"header-anchor",attrs:{href:"#why-this-parameter-consolidation"}},[e._v("#")]),e._v(" Why this parameter consolidation?")]),e._v(" "),a("p",[e._v("Reducing the number of parameters is good for UX. Fewer superfluous parameters makes\nrunning and operating a Tendermint network less confusing.")]),e._v(" "),a("p",[e._v("The Prevote and Precommit messages are both similar sizes, require similar amounts\nof processing so there is no strong need for them to be configured separately.")]),e._v(" "),a("p",[e._v("The "),a("code",[e._v("TimeoutPropose")]),e._v(" parameter governs how long Tendermint will wait for the proposed\nblock to be gossiped. Blocks are much larger than votes and therefore tend to be\ngossiped much more slowly. It therefore makes sense to keep "),a("code",[e._v("TimeoutPropose")]),e._v(" and\nthe "),a("code",[e._v("TimeoutProposeDelta")]),e._v(" as parameters separate from the vote timeouts.")]),e._v(" "),a("p",[a("code",[e._v("TimeoutCommit")]),e._v(" is used by chains to ensure that the network waits for the votes from\nslower validators before proceeding to the next height. Without this timeout, the votes\nfrom slower validators would consistently not be included in blocks and those validators\nwould not be counted as 'up' from the chain's perspective. Being down damages a validator's\nreputation and causes potential stakers to think twice before delegating to that validator.")]),e._v(" "),a("p",[a("code",[e._v("TimeoutCommit")]),e._v(" also prevents the network from producing the next height as soon as validators\non the fastest hardware with a summed voting power of +2/3 of the network's total have\ncompleted execution of the block. Allowing the network to proceed as soon as the fastest\n+2/3 completed execution would have a cumulative effect over heights, eventually\nleaving slower validators unable to participate in consensus at all. "),a("code",[e._v("TimeoutCommit")]),e._v("\ntherefore allows networks to have greater variability in hardware. Additional\ndiscussion of this can be found in "),a("a",{attrs:{href:"https://github.com/tendermint/tendermint/issues/5911#issuecomment-973560381",target:"_blank",rel:"noopener noreferrer"}},[e._v("tendermint issue 5911"),a("OutboundLink")],1),e._v("\nand "),a("a",{attrs:{href:"https://github.com/tendermint/spec/issues/359",target:"_blank",rel:"noopener noreferrer"}},[e._v("spec issue 359"),a("OutboundLink")],1),e._v(".")]),e._v(" "),a("h2",{attrs:{id:"alternative-approaches"}},[a("a",{staticClass:"header-anchor",attrs:{href:"#alternative-approaches"}},[e._v("#")]),e._v(" Alternative Approaches")]),e._v(" "),a("h3",{attrs:{id:"hardcode-the-parameters"}},[a("a",{staticClass:"header-anchor",attrs:{href:"#hardcode-the-parameters"}},[e._v("#")]),e._v(" Hardcode the parameters")]),e._v(" "),a("p",[e._v("Many Tendermint networks run on similar cloud-hosted infrastructure. Therefore,\nthey have similar bandwidth and machine resources. The timings for propagating votes\nand blocks are likely to be reasonably similar across networks. As a result, the\ntimeout parameters are good candidates for being hardcoded. Hardcoding the timeouts\nin Tendermint would mean entirely removing these parameters from any configuration\nthat could be altered by either an application or a node operator. Instead,\nTendermint would ship with a set of timeouts and all applications using Tendermint\nwould use this exact same set of values.")]),e._v(" "),a("p",[e._v("While Tendermint nodes often run with similar bandwidth and on similar cloud-hosted\nmachines, there are enough points of variability to make configuring\nconsensus timeouts meaningful. Namely, Tendermint network topologies are likely to be\nvery different from chain to chain. Additionally, applications may vary greatly in\nhow long the "),a("code",[e._v("Commit")]),e._v(" phase may take. Applications that perform more work during "),a("code",[e._v("Commit")]),e._v("\nrequire a longer "),a("code",[e._v("TimeoutCommit")]),e._v(" to allow the application to complete its work\nand be prepared for the next height.")]),e._v(" "),a("h2",{attrs:{id:"decision"}},[a("a",{staticClass:"header-anchor",attrs:{href:"#decision"}},[e._v("#")]),e._v(" Decision")]),e._v(" "),a("p",[e._v("The decision has been made to implement this work, with the caveat that the\nspecific mechanism for introducing the new parameters to chains is still ongoing.")]),e._v(" "),a("h2",{attrs:{id:"detailed-design"}},[a("a",{staticClass:"header-anchor",attrs:{href:"#detailed-design"}},[e._v("#")]),e._v(" Detailed Design")]),e._v(" "),a("h3",{attrs:{id:"new-consensus-parameters"}},[a("a",{staticClass:"header-anchor",attrs:{href:"#new-consensus-parameters"}},[e._v("#")]),e._v(" New Consensus Parameters")]),e._v(" "),a("p",[e._v("A new "),a("code",[e._v("TimeoutParams")]),e._v(" "),a("code",[e._v("message")]),e._v(" will be added to the [params.proto file][consensus-params-proto].\nThis message will have the following form:")]),e._v(" "),a("tm-code-block",{staticClass:"codeblock",attrs:{language:"proto",base64:"bWVzc2FnZSBUaW1lb3V0UGFyYW1zIHsKIGdvb2dsZS5wcm90b2J1Zi5EdXJhdGlvbiBwcm9wb3NlID0gMTsKIGdvb2dsZS5wcm90b2J1Zi5EdXJhdGlvbiBwcm9wb3NlX2RlbHRhID0gMjsKIGdvb2dsZS5wcm90b2J1Zi5EdXJhdGlvbiB2b3RlID0gMzsKIGdvb2dsZS5wcm90b2J1Zi5EdXJhdGlvbiB2b3RlX2RlbHRhID0gNDsKIGdvb2dsZS5wcm90b2J1Zi5EdXJhdGlvbiBjb21taXQgPSA1OwogYm9vbCBieXBhc3NfY29tbWl0X3RpbWVvdXQgPSA2Owp9Cg=="}}),e._v(" "),a("p",[e._v("This new message will be added as a field into the ["),a("code",[e._v("ConsensusParams")]),e._v("\nmessage][consensus-params-proto]. The same default values that are "),a("a",{attrs:{href:"https://github.com/tendermint/tendermint/blob/7cdf560173dee6773b80d1c574a06489d4c394fe/config/config.go#L955",target:"_blank",rel:"noopener noreferrer"}},[e._v("currently\nset for these parameters"),a("OutboundLink")],1),e._v(" in the local configuration\nfile will be used as the defaults for these new consensus parameters in the\n"),a("a",{attrs:{href:"https://github.com/tendermint/tendermint/blob/7cdf560173dee6773b80d1c574a06489d4c394fe/types/params.go#L79",target:"_blank",rel:"noopener noreferrer"}},[e._v("consensus parameter defaults"),a("OutboundLink")],1),e._v(".")]),e._v(" "),a("p",[e._v("The new consensus parameters will be subject to the same\n"),a("a",{attrs:{href:"https://github.com/tendermint/tendermint/blob/7cdf560173dee6773b80d1c574a06489d4c394fe/config/config.go#L1038",target:"_blank",rel:"noopener noreferrer"}},[e._v("validity rules"),a("OutboundLink")],1),e._v(" as the current configuration values,\nnamely, each value must be non-negative.")]),e._v(" "),a("h3",{attrs:{id:"migration"}},[a("a",{staticClass:"header-anchor",attrs:{href:"#migration"}},[e._v("#")]),e._v(" Migration")]),e._v(" "),a("p",[e._v("The new "),a("code",[e._v("ConsensusParameters")]),e._v(" will be added during an upcoming release. In this\nrelease, the old "),a("code",[e._v("config.toml")]),e._v(" parameters will cease to control the timeouts and\nan error will be logged on nodes that continue to specify these values. The specific\nmechanism by which these parameters will added to a chain is being discussed in\n"),a("a",{attrs:{href:"https://github.com/tendermint/tendermint/pull/7524",target:"_blank",rel:"noopener noreferrer"}},[e._v("RFC-009"),a("OutboundLink")],1),e._v(" and will be decided ahead of the next release.")]),e._v(" "),a("p",[e._v("The specific mechanism for adding these parameters depends on work related to\n"),a("a",{attrs:{href:"https://github.com/tendermint/spec/pull/222",target:"_blank",rel:"noopener noreferrer"}},[e._v("soft upgrades"),a("OutboundLink")],1),e._v(", which is still ongoing.")]),e._v(" "),a("h2",{attrs:{id:"consequences"}},[a("a",{staticClass:"header-anchor",attrs:{href:"#consequences"}},[e._v("#")]),e._v(" Consequences")]),e._v(" "),a("h3",{attrs:{id:"positive"}},[a("a",{staticClass:"header-anchor",attrs:{href:"#positive"}},[e._v("#")]),e._v(" Positive")]),e._v(" "),a("ul",[a("li",[e._v("Timeout parameters will be equal across all of the validators in a Tendermint network.")]),e._v(" "),a("li",[e._v("Remove superfluous timeout parameters.")])]),e._v(" "),a("h3",{attrs:{id:"negative"}},[a("a",{staticClass:"header-anchor",attrs:{href:"#negative"}},[e._v("#")]),e._v(" Negative")]),e._v(" "),a("h3",{attrs:{id:"neutral"}},[a("a",{staticClass:"header-anchor",attrs:{href:"#neutral"}},[e._v("#")]),e._v(" Neutral")]),e._v(" "),a("ul",[a("li",[e._v("Timeout parameters require consensus to change.")])]),e._v(" "),a("h2",{attrs:{id:"references"}},[a("a",{staticClass:"header-anchor",attrs:{href:"#references"}},[e._v("#")]),e._v(" References")])],1)}),[],!1,null,null,null);t.default=r.exports}}]);