import{N as e,O as l,Q as a,U as s,m as t,k as n,V as r,l as o,o as p,q as d,W as i,H as c,J as I,X as h,R as u,T as m,L as v,y as S}from"./preset_utils-B-rGdae_.chunk.js";import{a7 as g,aB as C,S as O,a8 as y,aC as A,aa as T,az as f,at as P,ab as k,ac as w,ad as L,ae as E,af as M,ag as b,ah as R,av as G,aj as K,ak as F,a6 as D,al as x,am as B,ao as W,P as q,ap as H,aw as j,aD as N,an as U,aq as V,a1 as z,ar as J,C as _,F as Q,R as X,T as Z}from"./detailed_results-DpSe3Rt6.chunk.js";e({fieldName:"latencyMs",label:"Latency",labelTooltip:"Player latency, in milliseconds. Adds a delay to actions that cannot be spell queued."}),l({fieldName:"assumeBleedActive",label:"Assume Bleed Always Active",labelTooltip:"Assume bleed always exists for 'Rend and Tear' calculations. Otherwise will only calculate based on own rip/rake/lacerate.",extraCssClasses:["within-raid-sim-hide"]});const Y={inputs:[a({fieldName:"minCombosForRip",label:"Min Rip CP",labelTooltip:"Combo Point threshold for allowing a Rip cast",float:!1,positive:!0}),a({fieldName:"maxWaitTime",label:"Max Wait Time",labelTooltip:"Max seconds to wait for an Energy tick to cast rather than powershifting",float:!0,positive:!0}),a({fieldName:"preroarDuration",label:"Pre-Roar Duration",labelTooltip:"Seconds remaining on a prior Savage Roar buff at the start of the pull",float:!0,positive:!0}),s({fieldName:"maintainFaerieFire",label:"Maintain Faerie Fire",labelTooltip:"If checked, bundle Faerie Fire refreshes with powershifts. Ignored if an external Faerie Fire debuff is selected in settings."}),s({fieldName:"precastTigersFury",label:"Pre-cast Tiger's Fury",labelTooltip:"If checked, model a pre-pull Tiger's Fury cast 3 seconds before initiating combat."}),s({fieldName:"useShredTrick",label:"Use Shred Trick",labelTooltip:'If checked, enable the "Shred trick" micro-optimization. This should only be used on short fight lengths with full powershifting uptime.'})]},$={type:"TypeAPL",prepullActions:[{action:{activateAura:{auraId:{spellId:768}}},doAtValue:{const:{val:"-10s"}}},{action:{activateAura:{auraId:{spellId:407988}}},doAtValue:{const:{val:"-8s"}}}],priorityList:[{action:{autocastOtherCooldowns:{}}},{action:{catOptimalRotationAction:{maxWaitTime:2,minCombosForRip:3}}}]},ee={type:"TypeAPL",prepullActions:[{action:{activateAura:{auraId:{spellId:768}}},doAtValue:{const:{val:"-10s"}}},{action:{activateAura:{auraId:{spellId:407988}}},doAtValue:{const:{val:"-8s"}}}],priorityList:[{action:{condition:{not:{val:{auraIsActive:{auraId:{spellId:768}}}}},castSpell:{spellId:{spellId:768}}}},{action:{condition:{and:{vals:[{cmp:{op:"OpLe",lhs:{currentManaPercent:{}},rhs:{const:{val:"40%"}}}},{cmp:{op:"OpLe",lhs:{currentEnergy:{}},rhs:{const:{val:"14"}}}}]}},castSpell:{spellId:{spellId:29166}}}},{action:{autocastOtherCooldowns:{}}},{action:{condition:{not:{val:{auraIsActive:{auraId:{spellId:407988}}}}},castSpell:{spellId:{spellId:407988}}}},{action:{condition:{cmp:{op:"OpLe",lhs:{currentEnergy:{}},rhs:{const:{val:"20"}}}},castSpell:{spellId:{spellId:417045}}}},{action:{condition:{not:{val:{auraIsActive:{sourceUnit:{type:"CurrentTarget"},auraId:{spellId:409828}}}}},castSpell:{spellId:{spellId:409828}}}},{action:{condition:{and:{vals:[{auraIsKnown:{auraId:{spellId:16870}}},{auraIsActive:{auraId:{spellId:16870}}}]}},castSpell:{spellId:{spellId:8992}}}},{action:{condition:{and:{vals:[{cmp:{op:"OpEq",lhs:{currentComboPoints:{}},rhs:{const:{val:"5"}}}},{cmp:{op:"OpGe",lhs:{auraRemainingTime:{auraId:{spellId:407988}}},rhs:{const:{val:"7"}}}}]}},castSpell:{spellId:{spellId:9493}}}},{action:{condition:{not:{val:{dotIsActive:{spellId:{spellId:414684}}}}},castSpell:{spellId:{spellId:414684}}}},{action:{condition:{not:{val:{auraIsActive:{sourceUnit:{type:"CurrentTarget"},auraId:{spellId:1823}}}}},castSpell:{spellId:{spellId:1823}}}},{action:{castSpell:{spellId:{spellId:8992}}}},{action:{condition:{and:{vals:[{cmp:{op:"OpLe",lhs:{currentEnergy:{}},rhs:{const:{val:"14"}}}},{cmp:{op:"OpGe",lhs:{currentMana:{}},rhs:{const:{val:"500"}}}},{auraIsActive:{auraId:{spellId:17061,rank:5}}}]}},castSpell:{spellId:{spellId:768}}}}]},le={type:"TypeAPL",prepullActions:[{action:{activateAura:{auraId:{spellId:768}}},doAtValue:{const:{val:"-10s"}}},{action:{activateAura:{auraId:{spellId:407988}}},doAtValue:{const:{val:"-8s"}}}],priorityList:[{action:{condition:{not:{val:{auraIsActive:{auraId:{spellId:768}}}}},castSpell:{spellId:{spellId:768}}}},{action:{condition:{and:{vals:[{cmp:{op:"OpLe",lhs:{math:{op:"OpAdd",lhs:{currentMana:{}},rhs:{const:{val:"1500.0"}}}},rhs:{math:{op:"OpDiv",lhs:{currentMana:{}},rhs:{currentManaPercent:{}}}}}},{cmp:{op:"OpLe",lhs:{currentEnergy:{}},rhs:{math:{op:"OpSub",lhs:{spellCurrentCost:{spellId:{spellId:1824,rank:3}}},rhs:{const:{val:"20.2"}}}}}},{cmp:{op:"OpGe",lhs:{currentMana:{}},rhs:{spellCurrentCost:{spellId:{spellId:768}}}}},{not:{val:{auraIsActive:{auraId:{spellId:417141}}}}},{auraIsKnown:{auraId:{spellId:17061,rank:5}}},{gcdIsReady:{}}]}},castSpell:{spellId:{itemId:12662}}}},{action:{condition:{and:{vals:[{cmp:{op:"OpLe",lhs:{currentMana:{}},rhs:{math:{op:"OpMul",lhs:{const:{val:"2.0"}},rhs:{spellCurrentCost:{spellId:{spellId:768}}}}}}},{cmp:{op:"OpLe",lhs:{currentEnergy:{}},rhs:{math:{op:"OpSub",lhs:{spellCurrentCost:{spellId:{spellId:1824,rank:3}}},rhs:{const:{val:"20.2"}}}}}},{cmp:{op:"OpGe",lhs:{currentMana:{}},rhs:{spellCurrentCost:{spellId:{spellId:768}}}}},{not:{val:{auraIsActive:{auraId:{spellId:417141}}}}},{auraIsKnown:{auraId:{spellId:17061,rank:5}}},{gcdIsReady:{}},{not:{val:{and:{vals:[{spellIsKnown:{spellId:{itemId:12662}}},{spellIsReady:{spellId:{itemId:12662}}}]}}}}]}},castSpell:{spellId:{otherId:"OtherActionPotion"}}}},{action:{condition:{and:{vals:[{cmp:{op:"OpLe",lhs:{currentManaPercent:{}},rhs:{const:{val:"40%"}}}},{cmp:{op:"OpLe",lhs:{currentEnergy:{}},rhs:{math:{op:"OpSub",lhs:{spellCurrentCost:{spellId:{spellId:1824,rank:3}}},rhs:{const:{val:"20.2"}}}}}},{cmp:{op:"OpGe",lhs:{currentMana:{}},rhs:{math:{op:"OpAdd",lhs:{spellCurrentCost:{spellId:{spellId:29166}}},rhs:{spellCurrentCost:{spellId:{spellId:768}}}}}}},{not:{val:{auraIsActive:{auraId:{spellId:417141}}}}},{auraIsKnown:{auraId:{spellId:17061,rank:5}}},{not:{val:{or:{vals:[{and:{vals:[{spellIsKnown:{spellId:{itemId:12662}}},{spellIsReady:{spellId:{itemId:12662}}}]}},{and:{vals:[{spellIsKnown:{spellId:{otherId:"OtherActionPotion"}}},{spellIsReady:{spellId:{otherId:"OtherActionPotion"}}}]}}]}}}}]}},castSpell:{spellId:{spellId:29166}}}},{action:{condition:{and:{vals:[{or:{vals:[{and:{vals:[{not:{val:{auraIsActive:{auraId:{spellId:417141}}}}},{cmp:{op:"OpLt",lhs:{currentEnergy:{}},rhs:{math:{op:"OpSub",lhs:{spellCurrentCost:{spellId:{spellId:1824,rank:3}}},rhs:{const:{val:"20.2"}}}}}}]}},{and:{vals:[{and:{vals:[{auraIsActive:{auraId:{spellId:417141}}},{spellIsKnown:{spellId:{spellId:417141}}}]}},{cmp:{op:"OpLt",lhs:{currentEnergy:{}},rhs:{spellCurrentCost:{spellId:{spellId:9829,rank:4}}}}}]}}]}},{cmp:{op:"OpGe",lhs:{currentMana:{}},rhs:{spellCurrentCost:{spellId:{spellId:768}}}}},{auraIsKnown:{auraId:{spellId:17061,rank:5}}}]}},castSpell:{spellId:{spellId:768}}}},{action:{autocastOtherCooldowns:{}}},{action:{condition:{or:{vals:[{cmp:{op:"OpLe",lhs:{currentEnergy:{}},rhs:{const:{val:"20"}}}},{and:{vals:[{cmp:{op:"OpGt",lhs:{timeToEnergyTick:{}},rhs:{gcdTimeToReady:{}}}},{cmp:{op:"OpLe",lhs:{currentEnergy:{}},rhs:{const:{val:"40"}}}}]}}]}},castSpell:{spellId:{spellId:417045}}}},{action:{condition:{not:{val:{auraIsActive:{auraId:{spellId:407988}}}}},castSpell:{spellId:{spellId:407988}}}},{action:{condition:{not:{val:{auraIsActive:{sourceUnit:{type:"CurrentTarget"},auraId:{spellId:409828}}}}},castSpell:{spellId:{spellId:409828}}}},{action:{condition:{and:{vals:[{auraIsKnown:{auraId:{spellId:16870}}},{auraIsActiveWithReactionTime:{auraId:{spellId:16870}}}]}},castSpell:{spellId:{spellId:9829}}}},{action:{condition:{and:{vals:[{cmp:{op:"OpEq",lhs:{currentComboPoints:{}},rhs:{const:{val:"5"}}}},{cmp:{op:"OpGe",lhs:{auraRemainingTime:{auraId:{spellId:407988}}},rhs:{const:{val:"7s"}}}},{cmp:{op:"OpGe",lhs:{remainingTime:{}},rhs:{const:{val:"10s"}}}},{not:{val:{dotIsActive:{spellId:{spellId:9752,rank:4}}}}}]}},castSpell:{spellId:{spellId:9752}}}},{action:{castSpell:{spellId:{spellId:9829}}}},{hide:!0,action:{condition:{and:{vals:[{runeIsEquipped:{runeId:{spellId:414684}}},{auraIsKnown:{auraId:{spellId:17061,rank:5}}},{cmp:{op:"OpGe",lhs:{currentMana:{}},rhs:{spellCurrentCost:{spellId:{spellId:768}}}}},{cmp:{op:"OpGt",lhs:{timeToEnergyTick:{}},rhs:{const:{val:"1.52s"}}}},{not:{val:{dotIsActive:{spellId:{spellId:414684}}}}}]}},castSpell:{spellId:{spellId:414684}}}},{action:{condition:{and:{vals:[{cmp:{op:"OpGe",lhs:{currentMana:{}},rhs:{spellCurrentCost:{spellId:{spellId:768}}}}},{auraIsKnown:{auraId:{spellId:17061,rank:5}}},{cmp:{op:"OpGt",lhs:{timeToEnergyTick:{}},rhs:{const:{val:"1.02s"}}}}]}},castSpell:{spellId:{spellId:409828}}}},{action:{condition:{and:{vals:[{cmp:{op:"OpGe",lhs:{currentMana:{}},rhs:{spellCurrentCost:{spellId:{spellId:768}}}}},{auraIsKnown:{auraId:{spellId:17061,rank:5}}},{cmp:{op:"OpGt",lhs:{timeToEnergyTick:{}},rhs:{const:{val:"1.02s"}}}},{not:{val:{dotIsActive:{spellId:{spellId:1824,rank:3}}}}},{not:{val:{runeIsEquipped:{runeId:{}}}}}]}},castSpell:{spellId:{spellId:1824,rank:3}}}}]},ae={type:"TypeAPL",prepullActions:[{action:{activateAura:{auraId:{spellId:407988}}},doAtValue:{const:{val:"-8s"}}}],priorityList:[{action:{condition:{not:{val:{auraIsActive:{auraId:{spellId:768}}}}},castSpell:{spellId:{spellId:768}}}},{action:{condition:{not:{val:{auraIsActive:{auraId:{spellId:407988}}}}},castSpell:{spellId:{spellId:407988}}}},{action:{condition:{and:{vals:[{cmp:{op:"OpLe",lhs:{math:{op:"OpAdd",lhs:{currentMana:{}},rhs:{const:{val:"1500.0"}}}},rhs:{math:{op:"OpDiv",lhs:{currentMana:{}},rhs:{currentManaPercent:{}}}}}},{cmp:{op:"OpLe",lhs:{currentEnergy:{}},rhs:{math:{op:"OpSub",lhs:{spellCurrentCost:{spellId:{spellId:9830,rank:5}}},rhs:{const:{val:"20.2"}}}}}},{cmp:{op:"OpGe",lhs:{currentMana:{}},rhs:{spellCurrentCost:{spellId:{spellId:768}}}}},{not:{val:{auraIsActive:{auraId:{spellId:417141}}}}}]}},castSpell:{spellId:{itemId:12662}}}},{action:{condition:{and:{vals:[{cmp:{op:"OpLe",lhs:{currentMana:{}},rhs:{math:{op:"OpMul",lhs:{const:{val:"2.0"}},rhs:{spellCurrentCost:{spellId:{spellId:768}}}}}}},{cmp:{op:"OpLe",lhs:{currentEnergy:{}},rhs:{math:{op:"OpSub",lhs:{spellCurrentCost:{spellId:{spellId:9830,rank:5}}},rhs:{const:{val:"20.2"}}}}}},{cmp:{op:"OpGe",lhs:{currentMana:{}},rhs:{spellCurrentCost:{spellId:{spellId:768}}}}},{not:{val:{auraIsActive:{auraId:{spellId:417141}}}}},{not:{val:{and:{vals:[{spellIsKnown:{spellId:{itemId:12662}}},{spellIsReady:{spellId:{itemId:12662}}}]}}}}]}},castSpell:{spellId:{otherId:"OtherActionPotion"}}}},{action:{condition:{and:{vals:[{cmp:{op:"OpLe",lhs:{currentManaPercent:{}},rhs:{const:{val:"40%"}}}},{cmp:{op:"OpLe",lhs:{currentEnergy:{}},rhs:{math:{op:"OpSub",lhs:{spellCurrentCost:{spellId:{spellId:9830,rank:5}}},rhs:{const:{val:"20.2"}}}}}},{cmp:{op:"OpGt",lhs:{currentMana:{}},rhs:{math:{op:"OpAdd",lhs:{spellCurrentCost:{spellId:{spellId:29166}}},rhs:{spellCurrentCost:{spellId:{spellId:768}}}}}}},{not:{val:{auraIsActive:{auraId:{spellId:417141}}}}},{not:{val:{or:{vals:[{and:{vals:[{spellIsKnown:{spellId:{itemId:12662}}},{spellIsReady:{spellId:{itemId:12662}}}]}},{and:{vals:[{spellIsKnown:{spellId:{otherId:"OtherActionPotion"}}},{spellIsReady:{spellId:{otherId:"OtherActionPotion"}}}]}}]}}}}]}},castSpell:{spellId:{spellId:29166}}}},{action:{condition:{and:{vals:[{auraIsKnown:{auraId:{spellId:17061,rank:5}}},{cmp:{op:"OpGe",lhs:{currentMana:{}},rhs:{spellCurrentCost:{spellId:{spellId:768}}}}},{or:{vals:[{and:{vals:[{not:{val:{auraIsActive:{auraId:{spellId:417141}}}}},{cmp:{op:"OpLt",lhs:{currentEnergy:{}},rhs:{math:{op:"OpSub",lhs:{spellCurrentCost:{spellId:{spellId:9830,rank:5}}},rhs:{const:{val:"20.2"}}}}}}]}},{and:{vals:[{runeIsEquipped:{runeId:{spellId:417141}}},{auraIsActive:{auraId:{spellId:417141}}},{cmp:{op:"OpLt",lhs:{currentEnergy:{}},rhs:{spellCurrentCost:{spellId:{spellId:9830,rank:5}}}}}]}}]}}]}},castSpell:{spellId:{spellId:768}}}},{action:{condition:{and:{vals:[{not:{val:{auraIsActive:{auraId:{spellId:417045}}}}},{or:{vals:[{not:{val:{energyThreshold:{threshold:-59}}}},{and:{vals:[{cmp:{op:"OpGt",lhs:{timeToEnergyTick:{}},rhs:{gcdTimeToReady:{}}}},{not:{val:{energyThreshold:{threshold:-39}}}}]}}]}}]}},castSpell:{spellId:{spellId:417045}}}},{action:{condition:{or:{vals:[{not:{val:{energyThreshold:{threshold:-79}}}},{and:{vals:[{cmp:{op:"OpGt",lhs:{timeToEnergyTick:{}},rhs:{gcdTimeToReady:{}}}},{not:{val:{energyThreshold:{threshold:-59}}}}]}}]}},castSpell:{spellId:{spellId:417045}}}},{action:{condition:{auraIsActive:{auraId:{spellId:417045}}},autocastOtherCooldowns:{}}},{action:{condition:{not:{val:{auraIsActive:{sourceUnit:{type:"CurrentTarget"},auraId:{spellId:409828}}}}},castSpell:{spellId:{spellId:409828}}}},{action:{condition:{and:{vals:[{cmp:{op:"OpGe",lhs:{numberTargets:{}},rhs:{const:{val:"4"}}}},{auraIsKnown:{auraId:{spellId:16870}}},{auraIsActiveWithReactionTime:{auraId:{spellId:16870}}}]}},castSpell:{spellId:{spellId:411128}}}},{action:{condition:{and:{vals:[{cmp:{op:"OpLt",lhs:{numberTargets:{}},rhs:{const:{val:"4"}}}},{auraIsKnown:{auraId:{spellId:16870}}},{auraIsActiveWithReactionTime:{auraId:{spellId:16870}}}]}},castSpell:{spellId:{spellId:9830,rank:5}}}},{action:{condition:{and:{vals:[{not:{val:{dotIsActive:{spellId:{spellId:9896,rank:6}}}}},{cmp:{op:"OpEq",lhs:{currentComboPoints:{}},rhs:{const:{val:"5"}}}},{cmp:{op:"OpGe",lhs:{remainingTime:{}},rhs:{const:{val:"10"}}}},{or:{vals:[{cmp:{op:"OpGe",lhs:{auraRemainingTime:{auraId:{spellId:407988}}},rhs:{const:{val:"8.0"}}}},{auraIsKnown:{auraId:{spellId:455873}}}]}}]}},castSpell:{spellId:{spellId:9896,rank:6}}}},{action:{condition:{not:{val:{dotIsActive:{spellId:{spellId:9904,rank:4}}}}},castSpell:{spellId:{spellId:9904,rank:4}}}},{action:{condition:{and:{vals:[{auraIsKnown:{auraId:{spellId:17061,rank:5}}},{cmp:{op:"OpGe",lhs:{currentMana:{}},rhs:{spellCurrentCost:{spellId:{spellId:768}}}}},{cmp:{op:"OpLt",lhs:{timeToEnergyTick:{}},rhs:{const:{val:"1.00"}}}},{cmp:{op:"OpGe",lhs:{math:{op:"OpAdd",lhs:{currentEnergy:{}},rhs:{const:{val:"20.2"}}}},rhs:{math:{op:"OpAdd",lhs:{spellCurrentCost:{spellId:{spellId:409828}}},rhs:{spellCurrentCost:{spellId:{spellId:409828}}}}}}},{cmp:{op:"OpLt",lhs:{math:{op:"OpAdd",lhs:{currentEnergy:{}},rhs:{const:{val:"20.2"}}}},rhs:{math:{op:"OpAdd",lhs:{spellCurrentCost:{spellId:{spellId:409828}}},rhs:{spellCurrentCost:{spellId:{spellId:9830,rank:5}}}}}}}]}},castSpell:{spellId:{spellId:409828}}}},{action:{condition:{and:{vals:[{auraIsKnown:{auraId:{spellId:455873}}},{cmp:{op:"OpEq",lhs:{currentComboPoints:{}},rhs:{const:{val:"5.0"}}}},{auraIsActive:{auraId:{spellId:407988}}},{not:{val:{auraIsActive:{auraId:{spellId:417141}}}}},{cmp:{op:"OpLe",lhs:{currentEnergy:{}},rhs:{const:{val:"53.0"}}}}]}},castSpell:{spellId:{spellId:31018}}}},{action:{condition:{cmp:{op:"OpGe",lhs:{numberTargets:{}},rhs:{const:{val:"4"}}}},castSpell:{spellId:{spellId:411128}}}},{action:{castSpell:{spellId:{spellId:9830,rank:5}}}},{action:{condition:{and:{vals:[{auraIsKnown:{auraId:{spellId:17061,rank:5}}},{cmp:{op:"OpGe",lhs:{currentMana:{}},rhs:{spellCurrentCost:{spellId:{spellId:768}}}}},{cmp:{op:"OpGt",lhs:{timeToEnergyTick:{}},rhs:{const:{val:"1.02"}}}}]}},castSpell:{spellId:{spellId:409828}}}}]},se={type:"TypeAPL",prepullActions:[{action:{activateAura:{auraId:{spellId:407988}}},doAtValue:{const:{val:"-8s"}},hide:!0}],priorityList:[{action:{condition:{not:{val:{auraIsActive:{auraId:{spellId:768}}}}},castSpell:{spellId:{spellId:768}}}},{action:{condition:{not:{val:{auraIsActive:{auraId:{spellId:407988}}}}},castSpell:{spellId:{spellId:407988}}}},{action:{condition:{and:{vals:[{cmp:{op:"OpLe",lhs:{math:{op:"OpAdd",lhs:{currentMana:{}},rhs:{const:{val:"1500.0"}}}},rhs:{math:{op:"OpDiv",lhs:{currentMana:{}},rhs:{currentManaPercent:{}}}}}},{cmp:{op:"OpLe",lhs:{currentEnergy:{}},rhs:{math:{op:"OpSub",lhs:{spellCurrentCost:{spellId:{spellId:9830,rank:5}}},rhs:{const:{val:"20.2"}}}}}},{cmp:{op:"OpGe",lhs:{currentMana:{}},rhs:{spellCurrentCost:{spellId:{spellId:768}}}}},{not:{val:{auraIsActive:{auraId:{spellId:417141}}}}},{auraIsKnown:{auraId:{spellId:17061,rank:5}}},{gcdIsReady:{}}]}},castSpell:{spellId:{itemId:12662}}}},{action:{condition:{and:{vals:[{cmp:{op:"OpLe",lhs:{currentMana:{}},rhs:{math:{op:"OpMul",lhs:{const:{val:"2.0"}},rhs:{spellCurrentCost:{spellId:{spellId:768}}}}}}},{cmp:{op:"OpLe",lhs:{currentEnergy:{}},rhs:{math:{op:"OpSub",lhs:{spellCurrentCost:{spellId:{spellId:9830,rank:5}}},rhs:{const:{val:"20.2"}}}}}},{cmp:{op:"OpGe",lhs:{currentMana:{}},rhs:{spellCurrentCost:{spellId:{spellId:768}}}}},{not:{val:{auraIsActive:{auraId:{spellId:417141}}}}},{auraIsKnown:{auraId:{spellId:17061,rank:5}}},{gcdIsReady:{}},{not:{val:{and:{vals:[{spellIsKnown:{spellId:{itemId:12662}}},{spellIsReady:{spellId:{itemId:12662}}}]}}}}]}},castSpell:{spellId:{otherId:"OtherActionPotion"}}}},{action:{condition:{and:{vals:[{cmp:{op:"OpLe",lhs:{currentManaPercent:{}},rhs:{const:{val:"40%"}}}},{cmp:{op:"OpLe",lhs:{currentEnergy:{}},rhs:{math:{op:"OpSub",lhs:{spellCurrentCost:{spellId:{spellId:9830,rank:5}}},rhs:{const:{val:"20.2"}}}}}},{cmp:{op:"OpGt",lhs:{currentMana:{}},rhs:{math:{op:"OpAdd",lhs:{spellCurrentCost:{spellId:{spellId:29166}}},rhs:{spellCurrentCost:{spellId:{spellId:768}}}}}}},{not:{val:{auraIsActive:{auraId:{spellId:417141}}}}},{auraIsKnown:{auraId:{spellId:17061,rank:5}}},{not:{val:{or:{vals:[{and:{vals:[{spellIsKnown:{spellId:{itemId:12662}}},{spellIsReady:{spellId:{itemId:12662}}}]}},{and:{vals:[{spellIsKnown:{spellId:{otherId:"OtherActionPotion"}}},{spellIsReady:{spellId:{otherId:"OtherActionPotion"}}}]}}]}}}}]}},castSpell:{spellId:{spellId:29166}}}},{action:{condition:{and:{vals:[{or:{vals:[{and:{vals:[{not:{val:{auraIsActive:{auraId:{spellId:417141}}}}},{cmp:{op:"OpLt",lhs:{currentEnergy:{}},rhs:{math:{op:"OpSub",lhs:{spellCurrentCost:{spellId:{spellId:9830,rank:5}}},rhs:{const:{val:"20.2"}}}}}}]}},{and:{vals:[{and:{vals:[{auraIsActive:{auraId:{spellId:417141}}},{spellIsKnown:{spellId:{spellId:417141}}}]}},{cmp:{op:"OpLt",lhs:{currentEnergy:{}},rhs:{spellCurrentCost:{spellId:{spellId:9830,rank:5}}}}}]}}]}},{cmp:{op:"OpGe",lhs:{currentMana:{}},rhs:{spellCurrentCost:{spellId:{spellId:768}}}}},{auraIsKnown:{auraId:{spellId:17061,rank:5}}}]}},castSpell:{spellId:{spellId:768}}}},{action:{condition:{and:{vals:[{not:{val:{auraIsActive:{auraId:{spellId:417045}}}}},{or:{vals:[{not:{val:{energyThreshold:{threshold:-59}}}},{and:{vals:[{cmp:{op:"OpGt",lhs:{timeToEnergyTick:{}},rhs:{gcdTimeToReady:{}}}},{not:{val:{energyThreshold:{threshold:-39}}}}]}}]}}]}},castSpell:{spellId:{spellId:417045}}}},{action:{condition:{or:{vals:[{not:{val:{energyThreshold:{threshold:-79}}}},{and:{vals:[{cmp:{op:"OpGt",lhs:{timeToEnergyTick:{}},rhs:{gcdTimeToReady:{}}}},{not:{val:{energyThreshold:{threshold:-59}}}}]}}]}},castSpell:{spellId:{spellId:417045}}}},{action:{condition:{auraIsActive:{auraId:{spellId:417045}}},autocastOtherCooldowns:{}}},{action:{condition:{not:{val:{auraIsActive:{sourceUnit:{type:"CurrentTarget"},auraId:{spellId:409828}}}}},castSpell:{spellId:{spellId:409828}}}},{action:{condition:{and:{vals:[{auraIsKnown:{auraId:{spellId:16870}}},{auraIsActiveWithReactionTime:{auraId:{spellId:16870}}}]}},castSpell:{spellId:{spellId:9830,rank:5}}}},{action:{condition:{and:{vals:[{cmp:{op:"OpEq",lhs:{currentComboPoints:{}},rhs:{const:{val:"5"}}}},{or:{vals:[{cmp:{op:"OpGe",lhs:{auraRemainingTime:{auraId:{spellId:407988}}},rhs:{const:{val:"8.0"}}}},{auraIsKnown:{auraId:{spellId:455873}}}]}},{cmp:{op:"OpGe",lhs:{remainingTime:{}},rhs:{const:{val:"10"}}}},{not:{val:{dotIsActive:{spellId:{spellId:9896,rank:6}}}}}]}},castSpell:{spellId:{spellId:9896,rank:6}}}},{action:{condition:{not:{val:{dotIsActive:{spellId:{spellId:9904,rank:4}}}}},castSpell:{spellId:{spellId:9904,rank:4}}}},{action:{condition:{and:{vals:[{cmp:{op:"OpGe",lhs:{currentMana:{}},rhs:{spellCurrentCost:{spellId:{spellId:768}}}}},{auraIsKnown:{auraId:{spellId:17061,rank:5}}},{cmp:{op:"OpLt",lhs:{timeToEnergyTick:{}},rhs:{const:{val:"1.00"}}}},{cmp:{op:"OpGe",lhs:{math:{op:"OpAdd",lhs:{currentEnergy:{}},rhs:{const:{val:"20.2"}}}},rhs:{math:{op:"OpAdd",lhs:{spellCurrentCost:{spellId:{spellId:409828}}},rhs:{spellCurrentCost:{spellId:{spellId:409828}}}}}}},{cmp:{op:"OpLt",lhs:{math:{op:"OpAdd",lhs:{currentEnergy:{}},rhs:{const:{val:"20.2"}}}},rhs:{math:{op:"OpAdd",lhs:{spellCurrentCost:{spellId:{spellId:409828}}},rhs:{spellCurrentCost:{spellId:{spellId:9830,rank:5}}}}}}}]}},castSpell:{spellId:{spellId:409828}}}},{action:{condition:{and:{vals:[{auraIsKnown:{auraId:{spellId:455873}}},{cmp:{op:"OpEq",lhs:{currentComboPoints:{}},rhs:{const:{val:"5.0"}}}},{auraIsActive:{auraId:{spellId:407988}}},{not:{val:{auraIsActive:{auraId:{spellId:417141}}}}},{cmp:{op:"OpLe",lhs:{currentEnergy:{}},rhs:{const:{val:"53.0"}}}}]}},castSpell:{spellId:{spellId:31018}}}},{action:{castSpell:{spellId:{spellId:9830,rank:5}}}},{action:{condition:{and:{vals:[{cmp:{op:"OpGe",lhs:{currentMana:{}},rhs:{spellCurrentCost:{spellId:{spellId:768}}}}},{auraIsKnown:{auraId:{spellId:17061,rank:5}}},{cmp:{op:"OpGt",lhs:{timeToEnergyTick:{}},rhs:{const:{val:"1.02"}}}}]}},castSpell:{spellId:{spellId:409828}}}}]},te={items:[{id:215166},{id:213344},{id:9647},{id:213307,enchant:849},{id:213313,enchant:866,rune:407977},{id:19590,enchant:856},{id:211423,enchant:856,rune:407995},{id:213322,rune:417141},{id:213332,rune:407988},{id:213341,enchant:849,rune:417046},{id:213284},{id:19512},{id:211449},{id:213348},{id:210741,enchant:34},{},{id:209576}]},ne={items:[{id:215166,enchant:7124,rune:417145},{id:13089},{id:220747,enchant:7328},{id:220615,enchant:849},{id:220779,enchant:928,rune:407977},{id:19590,enchant:856,rune:414719},{id:21319,enchant:1887,rune:407995},{id:213322,rune:417141},{id:220778,enchant:1508,rune:407988},{id:220780,enchant:1887,rune:417046},{id:19511},{id:12014,randomSuffix:692},{id:223195},{id:221307},{id:220596,enchant:34},{},{id:220606}]},re={items:[{id:226659,enchant:7124,rune:417145},{id:228685},{id:226665,enchant:7328},{id:228290,enchant:7564,rune:439510},{id:226661,enchant:1891,rune:407977},{id:226662,enchant:1885,rune:414719},{id:228257,enchant:927,rune:407995},{id:226660,rune:417141},{id:226666,enchant:1505,rune:407988},{id:226663,enchant:1887,rune:417046},{id:228286,rune:442896},{id:228261,rune:453622},{id:228078},{id:228089},{id:227683,enchant:1900},{},{id:22397}]},oe={items:[{id:231257,enchant:7124,rune:417145},{id:231803},{id:231259,enchant:2606},{id:230842,enchant:849,rune:439510},{id:231254,enchant:1891,rune:407977},{id:231261,enchant:1885,rune:414719},{id:232100,enchant:927,rune:407995},{id:232096,rune:417141},{id:232098,enchant:7615,rune:407988},{id:232101,enchant:1887,rune:417046},{id:230734,rune:453622},{id:228261,rune:442896},{id:231779},{id:230282},{id:224282,enchant:1900},{},{id:220606}]},pe=t("Phase 1",{items:[{id:211510},{id:209422},{id:209692},{id:213087,enchant:247},{id:211512,enchant:847,rune:407977},{id:209524,enchant:823},{id:211423,rune:407995},{id:209421},{id:10410,rune:407988},{id:211511,enchant:247},{id:20439},{id:6321},{id:211449},{id:4381},{id:209577,enchant:723},{},{id:209576}]},{customCondition:e=>25===e.getLevel()}),de=t("Phase 2",te,{customCondition:e=>40===e.getLevel()}),ie=t("Phase 3",ne,{customCondition:e=>50===e.getLevel()}),ce=t("Phase 4",re,{customCondition:e=>60===e.getLevel()}),Ie=t("Phase 5",oe,{customCondition:e=>60===e.getLevel()}),he={[g.Phase1]:[pe],[g.Phase2]:[de],[g.Phase3]:[ie],[g.Phase4]:[ce],[g.Phase5]:[Ie]},ue=he[g.Phase5][0],me=n("Phase 1",$,{customCondition:e=>25===e.getLevel()}),ve=n("Phase 2",ee,{customCondition:e=>40===e.getLevel()}),Se=n("Phase 3",le,{customCondition:e=>50===e.getLevel()}),ge=n("Phase 4",ae,{customCondition:e=>60===e.getLevel()}),Ce=n("Phase 5",se,{customCondition:e=>60===e.getLevel()}),Oe={[g.Phase1]:[me],[g.Phase2]:[ve],[g.Phase3]:[Se],[g.Phase4]:[ge],[g.Phase5]:[Ce]},ye={25:Oe[g.Phase1][0],40:Oe[g.Phase2][0],50:Oe[g.Phase3][0],60:Oe[g.Phase5][0]},Ae=C.create({maintainFaerieFire:!1,minCombosForRip:3,maxWaitTime:2,preroarDuration:26,precastTigersFury:!1,useShredTrick:!1});r("Simple Default",O.SpecFeralDruid,Ae);const Te=o("Level 25",y.create({talentsString:"500005001--05"}),{customCondition:e=>25===e.getLevel()}),fe=o("Level 40",y.create({talentsString:"-550002032320211-05"}),{customCondition:e=>40===e.getLevel()}),Pe=o("Level 50",y.create({talentsString:"500005301-5500020323002-05"}),{customCondition:e=>50===e.getLevel()}),ke=o("Level 50 LoTP",y.create({talentsString:"-5500020323202151-55"}),{customCondition:e=>50===e.getLevel()}),we=o("Level 60",y.create({talentsString:"500005301-5500020323202151-15"}),{customCondition:e=>60===e.getLevel()}),Le={[g.Phase1]:[Te],[g.Phase2]:[fe],[g.Phase3]:[Pe,ke],[g.Phase4]:[we],[g.Phase5]:[]},Ee=Le[g.Phase4][0],Me=A.create({latencyMs:100,assumeBleedActive:!0}),be=T.create({agilityElixir:f.ElixirOfTheHoneyBadger,attackPowerBuff:P.JujuMight,defaultConjured:k.ConjuredDemonicRune,defaultPotion:w.MajorManaPotion,dragonBreathChili:!0,enchantedSigil:L.WrathOfTheStormSigil,flask:E.FlaskOfDistilledWisdom,food:M.FoodSmokedDesertDumpling,mainHandImbue:b.ElementalSharpeningStone,manaRegenElixir:R.MagebloodPotion,mildlyIrradiatedRejuvPot:!0,miscConsumes:{catnip:!0,jujuEmber:!0},strengthBuff:G.JujuPower,zanzaBuff:K.ROIDS}),Re=F.create({arcaneBrilliance:!0,aspectOfTheLion:!0,battleShout:D.TristateEffectImproved,divineSpirit:!0,giftOfTheWild:D.TristateEffectImproved,graceOfAirTotem:D.TristateEffectImproved,leaderOfThePack:!0,manaSpringTotem:D.TristateEffectRegular,strengthOfEarthTotem:D.TristateEffectImproved}),Ge=x.create({blessingOfKings:!0,blessingOfMight:D.TristateEffectImproved,blessingOfWisdom:D.TristateEffectImproved,fengusFerocity:!0,mightOfStormwind:!0,rallyingCryOfTheDragonslayer:!0,saygesFortune:B.SaygesDamage,songflowerSerenade:!0,spiritOfZandalar:!0,valorOfAzeroth:!0,warchiefsBlessing:!0}),Ke=W.create({curseOfRecklessness:!0,exposeArmor:D.TristateEffectImproved,faerieFire:!0,homunculi:70}),Fe={profession1:q.Enchanting,profession2:q.Alchemy},De=p(O.SpecFeralDruid,{cssClass:"feral-druid-sim-ui",cssScheme:"druid",knownIssues:[],warnings:[],epStats:[H.StatMana,H.StatStrength,H.StatAgility,H.StatIntellect,H.StatSpirit,H.StatAttackPower,H.StatFeralAttackPower,H.StatMeleeHit,H.StatMeleeCrit,H.StatMeleeHaste,H.StatExpertise,H.StatMP5],epPseudoStats:[j.PseudoStatBonusPhysicalDamage],epReferenceStat:H.StatAttackPower,displayStats:[H.StatMana,H.StatStrength,H.StatAgility,H.StatIntellect,H.StatSpirit,H.StatAttackPower,H.StatFeralAttackPower,H.StatMeleeHit,H.StatMeleeCrit,H.StatExpertise,H.StatMP5],displayPseudoStats:[j.PseudoStatBonusPhysicalDamage],defaults:{gear:ue.gear,epWeights:d.fromMap({[H.StatStrength]:2.38,[H.StatAgility]:2.35,[H.StatAttackPower]:1,[H.StatFeralAttackPower]:1,[H.StatMeleeHit]:24.46,[H.StatMeleeCrit]:16.67,[H.StatMana]:.04,[H.StatIntellect]:.67,[H.StatSpirit]:.08,[H.StatMP5]:.46,[H.StatFireResistance]:.5},{}),consumes:be,rotationType:N.TypeSimple,simpleRotation:Ae,talents:Ee.data,specOptions:Me,other:Fe,raidBuffs:Re,partyBuffs:U.create({}),individualBuffs:Ge,debuffs:Ke},playerIconInputs:[],rotationInputs:Y,includeBuffDebuffInputs:[i,c,I,h],excludeBuffDebuffInputs:[b.ElementalSharpeningStone,b.DenseSharpeningStone,b.WildStrikes],otherInputs:{inputs:[u,m,v]},itemSwapConfig:{itemSlots:[V.ItemSlotMainHand,V.ItemSlotOffHand,V.ItemSlotRanged]},encounterPicker:{showExecuteProportion:!1},presets:{talents:[...Le[g.Phase5],...Le[g.Phase4],...Le[g.Phase3],...Le[g.Phase2],...Le[g.Phase1]],rotations:[...Oe[g.Phase4],...Oe[g.Phase3],...Oe[g.Phase2],...Oe[g.Phase1]],gear:[...he[g.Phase5],...he[g.Phase4],...he[g.Phase3],...he[g.Phase2],...he[g.Phase1]]},autoRotation:e=>ye[e.getLevel()].rotation.rotation,raidSimPresets:[{spec:O.SpecFeralDruid,tooltip:z[O.SpecFeralDruid],defaultName:"Cat",iconUrl:J(_.ClassDruid,3),talents:Ee.data,specOptions:Me,consumes:be,defaultFactionRaces:{[Q.Unknown]:X.RaceUnknown,[Q.Alliance]:X.RaceNightElf,[Q.Horde]:X.RaceTauren},defaultGear:{[Q.Unknown]:{},[Q.Alliance]:{1:he[g.Phase1][0].gear},[Q.Horde]:{1:he[g.Phase1][0].gear}}}]});class xe extends S{constructor(e,l){super(e,l,De)}calcArpTarget(e){let l=1399;e.hasTrinket(45931)?l-=751:e.hasTrinket(40256)&&(l-=612);const a=e.getEquippedItem(V.ItemSlotMainHand);return null!=a&&null!=a.enchant&&3225==a.enchant.effectId&&(l-=120),l}calcCritCap(e){let l=0;return e.hasRelic(47668)&&(l+=200),e.hasRelic(50456)&&(l+=220),(e.hasTrinket(47131)||e.hasTrinket(47464))&&(l+=510),(e.hasTrinket(47115)||e.hasTrinket(47303))&&(l+=450),(e.hasTrinket(44253)||e.hasTrinket(42987))&&(l+=300),(new d).withStat(H.StatMeleeCrit,45.91*(77.8-1.1*l*1.06*1.02/83.33))}async updateGear(e){return this.player.setGear(Z.nextEventID(),e),await this.sim.updateCharacterStats(Z.nextEventID()),d.fromProto(this.player.getCurrentStats().finalStats)}}export{xe as F};