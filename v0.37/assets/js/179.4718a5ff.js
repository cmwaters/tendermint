(window.webpackJsonp=window.webpackJsonp||[]).push([[179],{791:function(e,t,o){"use strict";o.r(t);var r=o(1),n=Object(r.a)({},(function(){var e=this,t=e.$createElement,o=e._self._c||t;return o("ContentSlotsDistributor",{attrs:{"slot-key":e.$parent.slotKey}},[o("h1",{attrs:{id:"mempool"}},[o("a",{staticClass:"header-anchor",attrs:{href:"#mempool"}},[e._v("#")]),e._v(" Mempool")]),e._v(" "),o("h2",{attrs:{id:"transaction-ordering"}},[o("a",{staticClass:"header-anchor",attrs:{href:"#transaction-ordering"}},[e._v("#")]),e._v(" Transaction ordering")]),e._v(" "),o("p",[e._v("Currently, there's no ordering of transactions other than the order they've\narrived (via RPC or from other nodes).")]),e._v(" "),o("p",[e._v("So the only way to specify the order is to send them to a single node.")]),e._v(" "),o("p",[e._v("valA:")]),e._v(" "),o("ul",[o("li",[o("code",[e._v("tx1")])]),e._v(" "),o("li",[o("code",[e._v("tx2")])]),e._v(" "),o("li",[o("code",[e._v("tx3")])])]),e._v(" "),o("p",[e._v("If the transactions are split up across different nodes, there's no way to\nensure they are processed in the expected order.")]),e._v(" "),o("p",[e._v("valA:")]),e._v(" "),o("ul",[o("li",[o("code",[e._v("tx1")])]),e._v(" "),o("li",[o("code",[e._v("tx2")])])]),e._v(" "),o("p",[e._v("valB:")]),e._v(" "),o("ul",[o("li",[o("code",[e._v("tx3")])])]),e._v(" "),o("p",[e._v("If valB is the proposer, the order might be:")]),e._v(" "),o("ul",[o("li",[o("code",[e._v("tx3")])]),e._v(" "),o("li",[o("code",[e._v("tx1")])]),e._v(" "),o("li",[o("code",[e._v("tx2")])])]),e._v(" "),o("p",[e._v("If valA is the proposer, the order might be:")]),e._v(" "),o("ul",[o("li",[o("code",[e._v("tx1")])]),e._v(" "),o("li",[o("code",[e._v("tx2")])]),e._v(" "),o("li",[o("code",[e._v("tx3")])])]),e._v(" "),o("p",[e._v("That said, if the transactions contain some internal value, like an\norder/nonce/sequence number, the application can reject transactions that are\nout of order. So if a node receives "),o("code",[e._v("tx3")]),e._v(", then "),o("code",[e._v("tx1")]),e._v(", it can reject "),o("code",[e._v("tx3")]),e._v(" and then\naccept "),o("code",[e._v("tx1")]),e._v(". The sender can then retry sending "),o("code",[e._v("tx3")]),e._v(", which should probably be\nrejected until the node has seen "),o("code",[e._v("tx2")]),e._v(".")])])}),[],!1,null,null,null);t.default=n.exports}}]);