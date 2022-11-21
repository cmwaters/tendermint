(window.webpackJsonp=window.webpackJsonp||[]).push([[30],{622:function(e,t,o){e.exports=o.p+"assets/img/sentry_layout.4c8ab19d.png"},623:function(e,t,o){e.exports=o.p+"assets/img/local_config.0b47aff9.png"},797:function(e,t,o){"use strict";o.r(t);var n=o(1),a=Object(n.a)({},(function(){var e=this,t=e.$createElement,n=e._self._c||t;return n("ContentSlotsDistributor",{attrs:{"slot-key":e.$parent.slotKey}},[n("h1",{attrs:{id:"validators"}},[n("a",{staticClass:"header-anchor",attrs:{href:"#validators"}},[e._v("#")]),e._v(" Validators")]),e._v(" "),n("p",[e._v("Validators are responsible for committing new blocks in the blockchain.\nThese validators participate in the consensus protocol by broadcasting\n"),n("em",[e._v("votes")]),e._v(" which contain cryptographic signatures signed by each\nvalidator's private key.")]),e._v(" "),n("p",[e._v('Some Proof-of-Stake consensus algorithms aim to create a "completely"\ndecentralized system where all stakeholders (even those who are not\nalways available online) participate in the committing of blocks.\nTendermint has a different approach to block creation. Validators are\nexpected to be online, and the set of validators is permissioned/curated\nby some external process. Proof-of-stake is not required, but can be\nimplemented on top of Tendermint consensus. That is, validators may be\nrequired to post collateral on-chain, off-chain, or may not be required\nto post any collateral at all.')]),e._v(" "),n("p",[e._v('Validators have a cryptographic key-pair and an associated amount of\n"voting power". Voting power need not be the same.')]),e._v(" "),n("h2",{attrs:{id:"becoming-a-validator"}},[n("a",{staticClass:"header-anchor",attrs:{href:"#becoming-a-validator"}},[e._v("#")]),e._v(" Becoming a Validator")]),e._v(" "),n("p",[e._v("There are two ways to become validator.")]),e._v(" "),n("ol",[n("li",[e._v("They can be pre-established in the "),n("RouterLink",{attrs:{to:"/tendermint-core/using-tendermint.html#genesis"}},[e._v("genesis state")])],1),e._v(" "),n("li",[e._v("The ABCI app responds to the EndBlock message with changes to the\nexisting validator set.")])]),e._v(" "),n("h2",{attrs:{id:"setting-up-a-validator"}},[n("a",{staticClass:"header-anchor",attrs:{href:"#setting-up-a-validator"}},[e._v("#")]),e._v(" Setting up a Validator")]),e._v(" "),n("p",[e._v("When setting up a validator there are countless ways to configure your setup. This guide is aimed at showing one of them, the sentry node design. This design is mainly for DDOS prevention.")]),e._v(" "),n("h3",{attrs:{id:"network-layout"}},[n("a",{staticClass:"header-anchor",attrs:{href:"#network-layout"}},[e._v("#")]),e._v(" Network Layout")]),e._v(" "),n("p",[n("img",{attrs:{src:o(622),alt:"ALT Network Layout"}})]),e._v(" "),n("p",[e._v("The diagram is based on AWS, other cloud providers will have similar solutions to design a solution. Running nodes is not limited to cloud providers, you can run nodes on bare metal systems as well. The architecture will be the same no matter which setup you decide to go with.")]),e._v(" "),n("p",[e._v("The proposed network diagram is similar to the classical backend/frontend separation of services in a corporate environment. The “backend” in this case is the private network of the validator in the data center. The data center network might involve multiple subnets, firewalls and redundancy devices, which is not detailed on this diagram. The important point is that the data center allows direct connectivity to the chosen cloud environment. Amazon AWS has “Direct Connect”, while Google Cloud has “Partner Interconnect”. This is a dedicated connection to the cloud provider (usually directly to your virtual private cloud instance in one of the regions).")]),e._v(" "),n("p",[e._v("All sentry nodes (the “frontend”) connect to the validator using this private connection. The validator does not have a public IP address to provide its services.")]),e._v(" "),n("p",[e._v("Amazon has multiple availability zones within a region. One can install sentry nodes in other regions too. In this case the second, third and further regions need to have a private connection to the validator node. This can be achieved by VPC Peering (“VPC Network Peering” in Google Cloud). In this case, the second, third and further region sentry nodes will be directed to the first region and through the direct connect to the data center, arriving to the validator.")]),e._v(" "),n("p",[e._v("A more persistent solution (not detailed on the diagram) is to have multiple direct connections to different regions from the data center. This way VPC Peering is not mandatory, although still beneficial for the sentry nodes. This overcomes the risk of depending on one region. It is more costly.")]),e._v(" "),n("h3",{attrs:{id:"local-configuration"}},[n("a",{staticClass:"header-anchor",attrs:{href:"#local-configuration"}},[e._v("#")]),e._v(" Local Configuration")]),e._v(" "),n("p",[n("img",{attrs:{src:o(623),alt:"ALT Local Configuration"}})]),e._v(" "),n("p",[e._v("The validator will only talk to the sentry that are provided, the sentry nodes will communicate to the validator via a secret connection and the rest of the network through a normal connection. The sentry nodes do have the option of communicating with each other as well.")]),e._v(" "),n("p",[e._v("When initializing nodes there are five parameters in the "),n("code",[e._v("config.toml")]),e._v(" that may need to be altered.")]),e._v(" "),n("ul",[n("li",[n("code",[e._v("pex:")]),e._v(" boolean. This turns the peer exchange reactor on or off for a node. When "),n("code",[e._v("pex=false")]),e._v(", only the "),n("code",[e._v("persistent_peers")]),e._v(" list is available for connection.")]),e._v(" "),n("li",[n("code",[e._v("persistent_peers:")]),e._v(" a comma separated list of "),n("code",[e._v("nodeID@ip:port")]),e._v(" values that define a list of peers that are expected to be online at all times. This is necessary at first startup because by setting "),n("code",[e._v("pex=false")]),e._v(" the node will not be able to join the network.")]),e._v(" "),n("li",[n("code",[e._v("unconditional_peer_ids:")]),e._v(" comma separated list of nodeID's. These nodes will be connected to no matter the limits of inbound and outbound peers. This is useful for when sentry nodes have full address books.")]),e._v(" "),n("li",[n("code",[e._v("private_peer_ids:")]),e._v(" comma separated list of nodeID's. These nodes will not be gossiped to the network. This is an important field as you do not want your validator IP gossiped to the network.")]),e._v(" "),n("li",[n("code",[e._v("addr_book_strict:")]),e._v(" boolean. By default nodes with a routable address will be considered for connection. If this setting is turned off (false), non-routable IP addresses, like addresses in a private network can be added to the address book.")]),e._v(" "),n("li",[n("code",[e._v("double_sign_check_height")]),e._v(" int64 height.  How many blocks to look back to check existence of the node's consensus votes before joining consensus When non-zero, the node will panic upon restart if the same consensus key was used to sign {double_sign_check_height} last blocks. So, validators should stop the state machine, wait for some blocks, and then restart the state machine to avoid panic.")])]),e._v(" "),n("h4",{attrs:{id:"validator-node-configuration"}},[n("a",{staticClass:"header-anchor",attrs:{href:"#validator-node-configuration"}},[e._v("#")]),e._v(" Validator Node Configuration")]),e._v(" "),n("table",[n("thead",[n("tr",[n("th",[e._v("Config Option")]),e._v(" "),n("th",[e._v("Setting")])])]),e._v(" "),n("tbody",[n("tr",[n("td",[e._v("pex")]),e._v(" "),n("td",[e._v("false")])]),e._v(" "),n("tr",[n("td",[e._v("persistent_peers")]),e._v(" "),n("td",[e._v("list of sentry nodes")])]),e._v(" "),n("tr",[n("td",[e._v("private_peer_ids")]),e._v(" "),n("td",[e._v("none")])]),e._v(" "),n("tr",[n("td",[e._v("unconditional_peer_ids")]),e._v(" "),n("td",[e._v("optionally sentry node IDs")])]),e._v(" "),n("tr",[n("td",[e._v("addr_book_strict")]),e._v(" "),n("td",[e._v("false")])]),e._v(" "),n("tr",[n("td",[e._v("double_sign_check_height")]),e._v(" "),n("td",[e._v("10")])])])]),e._v(" "),n("p",[e._v("The validator node should have "),n("code",[e._v("pex=false")]),e._v(" so it does not gossip to the entire network. The persistent peers will be your sentry nodes. Private peers can be left empty as the validator is not trying to hide who it is communicating with. Setting unconditional peers is optional for a validator because they will not have a full address books.")]),e._v(" "),n("h4",{attrs:{id:"sentry-node-configuration"}},[n("a",{staticClass:"header-anchor",attrs:{href:"#sentry-node-configuration"}},[e._v("#")]),e._v(" Sentry Node Configuration")]),e._v(" "),n("table",[n("thead",[n("tr",[n("th",[e._v("Config Option")]),e._v(" "),n("th",[e._v("Setting")])])]),e._v(" "),n("tbody",[n("tr",[n("td",[e._v("pex")]),e._v(" "),n("td",[e._v("true")])]),e._v(" "),n("tr",[n("td",[e._v("persistent_peers")]),e._v(" "),n("td",[e._v("validator node, optionally other sentry nodes")])]),e._v(" "),n("tr",[n("td",[e._v("private_peer_ids")]),e._v(" "),n("td",[e._v("validator node ID")])]),e._v(" "),n("tr",[n("td",[e._v("unconditional_peer_ids")]),e._v(" "),n("td",[e._v("validator node ID, optionally sentry node IDs")])]),e._v(" "),n("tr",[n("td",[e._v("addr_book_strict")]),e._v(" "),n("td",[e._v("false")])])])]),e._v(" "),n("p",[e._v("The sentry nodes should be able to talk to the entire network hence why "),n("code",[e._v("pex=true")]),e._v(". The persistent peers of a sentry node will be the validator, and optionally other sentry nodes. The sentry nodes should make sure that they do not gossip the validator's ip, to do this you must put the validators nodeID as a private peer. The unconditional peer IDs will be the validator ID and optionally other sentry nodes.")]),e._v(" "),n("blockquote",[n("p",[e._v("Note: Do not forget to secure your node's firewalls when setting them up.")])]),e._v(" "),n("p",[e._v("More Information can be found at these links:")]),e._v(" "),n("ul",[n("li",[n("a",{attrs:{href:"https://kb.certus.one/",target:"_blank",rel:"noopener noreferrer"}},[e._v("https://kb.certus.one/"),n("OutboundLink")],1)]),e._v(" "),n("li",[n("a",{attrs:{href:"https://forum.cosmos.network/t/sentry-node-architecture-overview/454",target:"_blank",rel:"noopener noreferrer"}},[e._v("https://forum.cosmos.network/t/sentry-node-architecture-overview/454"),n("OutboundLink")],1)])]),e._v(" "),n("h3",{attrs:{id:"validator-keys"}},[n("a",{staticClass:"header-anchor",attrs:{href:"#validator-keys"}},[e._v("#")]),e._v(" Validator keys")]),e._v(" "),n("p",[e._v("Protecting a validator's consensus key is the most important factor to take in when designing your setup. The key that a validator is given upon creation of the node is called a consensus key, it has to be online at all times in order to vote on blocks. It is "),n("strong",[e._v("not recommended")]),e._v(" to merely hold your private key in the default json file ("),n("code",[e._v("priv_validator_key.json")]),e._v("). Fortunately, the "),n("a",{attrs:{href:"https://interchain.io/",target:"_blank",rel:"noopener noreferrer"}},[e._v("Interchain Foundation"),n("OutboundLink")],1),e._v(" has worked with a team to build a key management server for validators. You can find documentation on how to use it "),n("a",{attrs:{href:"https://github.com/iqlusioninc/tmkms",target:"_blank",rel:"noopener noreferrer"}},[e._v("here"),n("OutboundLink")],1),e._v(", it is used extensively in production. You are not limited to using this tool, there are also "),n("a",{attrs:{href:"https://safenet.gemalto.com/data-encryption/hardware-security-modules-hsms/",target:"_blank",rel:"noopener noreferrer"}},[e._v("HSMs"),n("OutboundLink")],1),e._v(", there is not a recommended HSM.")]),e._v(" "),n("p",[e._v("Currently Tendermint uses "),n("a",{attrs:{href:"https://ed25519.cr.yp.to/",target:"_blank",rel:"noopener noreferrer"}},[e._v("Ed25519"),n("OutboundLink")],1),e._v(" keys which are widely supported across the security sector and HSMs.")]),e._v(" "),n("h2",{attrs:{id:"committing-a-block"}},[n("a",{staticClass:"header-anchor",attrs:{href:"#committing-a-block"}},[e._v("#")]),e._v(" Committing a Block")]),e._v(" "),n("blockquote",[n("p",[n("strong",[e._v('+2/3 is short for "more than 2/3"')])])]),e._v(" "),n("p",[e._v("A block is committed when +2/3 of the validator set sign "),n("a",{attrs:{href:"https://github.com/tendermint/spec/blob/953523c3cb99fdb8c8f7a2d21e3a99094279e9de/spec/blockchain/blockchain.md#vote",target:"_blank",rel:"noopener noreferrer"}},[e._v("precommit\nvotes"),n("OutboundLink")],1),e._v(" for that block at the same "),n("code",[e._v("round")]),e._v(".\nThe +2/3 set of precommit votes is called a\n"),n("a",{attrs:{href:"https://github.com/tendermint/spec/blob/953523c3cb99fdb8c8f7a2d21e3a99094279e9de/spec/blockchain/blockchain.md#commit",target:"_blank",rel:"noopener noreferrer"}},[n("em",[e._v("commit")]),n("OutboundLink")],1),e._v(". While any +2/3 set of\nprecommits for the same block at the same height&round can serve as\nvalidation, the canonical commit is included in the next block (see\n"),n("a",{attrs:{href:"https://github.com/tendermint/spec/blob/953523c3cb99fdb8c8f7a2d21e3a99094279e9de/spec/blockchain/blockchain.md#lastcommit",target:"_blank",rel:"noopener noreferrer"}},[e._v("LastCommit"),n("OutboundLink")],1),e._v(").")])])}),[],!1,null,null,null);t.default=a.exports}}]);