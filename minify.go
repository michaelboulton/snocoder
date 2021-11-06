package main

import (
	"bytes"
	"compress/flate"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var minifyCmd = &cobra.Command{
	Use: "minify",
	RunE: func(cmd *cobra.Command, args []string) error {
		return minify()
	},
}

const dict = `7g5nnjrcwchhc2227g5nnmp4q4hhc2257g5nn384lxbsk2277g5nn93k7at2e2297g5nnjrc98hhc2297g5nml2wih9ac2277g5nmogcfk3e422b7g5nmvqovs49i2267g5nn7ko6ohhc22f7g5nmc9owjpie22e7g5nnq6gwk3e422g7g5nmhlkaqt2e22h7g5nnaigfrimu22n7g5nnnocofimu22s7g5nmpv8pqt2e22x7g5nngdwldbsk2357g5nmhlkn5bsk2377g5nn28s7g49i22j7g5nmubsatbsk23a7g5nmf7g9wraq23b7g5nnqm4ix9ac2347g5nncgwhbimu23h7g5nmpv8tchhc23j7g5nn8jw2khhc23p7g5nn384nit2e23h7g5nn3845et2e23k7g5nnmp4h43e423l7g5nnn8oqg49i2347g5nno7welbsk23t7g5nmjk4dchhc23v7g5nmwacpo49i2377g5nmg6o92t2e23r7g5nnmp4cjrwi23u7g5nmre4vghhc2427g5nmqeseg3e423s7g5nn7ko7mt2e23t7g5nmubsqt9ac23u7g5nmre48rpie23x7g5nnlpsarrwi2437g5nmovwppbsk24d7g5nmiks88raq2457g5nno7wc4raq2467g5nnq6g9bpie2467g5nnb24avimu24e7g5nmsdccghhc24g7g5nn3844shhc24h7g5nnicgqpbsk24k7g5nmnwoacraq2497g5nn93khgraq24c7g5nn65sn4hhc24p7g5nn84c6o3e424d7g5nmmhsup9ac24j7g5nnnoc63pie24i7g5nnskkxbpie24j7g5nmc9o6o3e424r7g5nmens6t9ac24u7g5nnn8o9oraq2527g5nmurgl83e424u7g5nmdok3t9ac2527g5nn754enpie24q7g5nn56kgchhc25c7g5nmh5wm3rwi25d7g5nmenshhbsk25g7g5nn7koc3rwi25f7g5nncgwnchhc25q7g5nmpv8p7pie24x7g5nna2stmt2e2597g5nmd8wmh9ac25i7g5nnkqk74raq2677g5nn2ogrd9ac25v7g5nnb24hfimu2627g5nnhd8t59ac26c7g5nmd8whw3e42657g5nml2wpghhc26w7g5nn3844pbsk2727g5nnq6gtlbsk2737g5nnhssvwraq26t7g5nnb24r99ac26j7g5nnr5o8ut2e26b7g5nmm28cfpie26c7g5nmnh4wnrwi26h7g5nn65sx849i25p7g5nmubsbqt2e26g7g5nmsdc9k3e426o7g5nmwacbjimu26w7g5nn93kc59ac26s7g5nmsdccchhc27j7g5nmx9kebpie26j7g5nnonk3graq27g7g5nmd8wt849i2657g5nnicg6ut2e26w7g5nmts85g3e42777g5nn2ogget2e27b7g5nmfn4het2e27i7g5nmvqoxbimu27s7g5nmdokvg3e427u7g5nn5m8n9bsk28r7g5nmx9kupbsk28s7g5nmens53imu2897g5nns4wenimu28c7g5nmogc2t9ac2887g5nn56kptbsk2947g5nmjk4eoraq28k7g5nn84cwt9ac28b7g5nnmp4w2t2e2877g5nmh5whwraq28p7g5nmh5wcwhhc29a7g5nnlpsput2e28b7g5nno7whoraq28v7g5nmrtowjimu2927g5nns4wbvrwi28e7g5nn65sjghhc29l7g5nnonklmt2e28k7g5nnkawxs3e428o7g5nnla8u3rwi28o7g5nnrlcvohhc29q7g5nncwkc2t2e28r7g5nmubs48raq29m7g5nnev4kut2e28w7g5nnq6gxg49i27u7g5nmovw9vimu29p7g5nmd8wml9ac29h7g5nnjrc2mt2e2947g5nnis478hhc2aa7g5nncgwhxbsk2a77g5nnefg3l9ac29s7g5nnpmssw3e429f7g5nnm9gbw49i28c7g5nmlikbt9ac29x7g5nnjrc8tbsk2ag7g5nn5m8fat2e29m7g5nngtkrrrwi29m7g5nnlpsjshhc2am7g5nmh5wjk49i28p7g5nnmp45jrwi29w7g5nn47chjpie29t7g5nmpv8bwhhc2av7g5nml2wn8hhc2b67g5nnla8vnpie2a47g5nmdokepbsk2b97g5nmnwo8jrwi2ag7g5nn384jl9ac2au7g5nmi584ut2e2ad7g5nmrtop3rwi2aq7g5nmk3ognrwi2au7g5nmvb4dbpie2ae7g5nmvqo7c3e42ai7g5nnm9gnhbsk2bm7g5nn56kidbsk2bo7g5nnhd84k49i29n7g5nmvqovw49i29o7g5nn56kwx9ac2bb7g5nmxp8edbsk2bv7g5nnbhokshhc2bo7g5nmens97pie2an7g5nnkqkis3e42ar7g5nmcpcvnimu2bi7g5nno7wrpbsk2cd7g5nmd8wbwraq2c67g5nn8jw4vpie2b27g5nncwkl59ac2bs7g5nmre4i6t2e2bh7g5nnicgbqt2e2bl7g5nn6lg883e42bg7g5nmmxgad9ac2c27g5nmnh4rrrwi2bu7g5nmwpwohbsk2cx7g5nmg6o3g3e42bp7g5nmx9kakhhc2ch7g5nndg8bt9ac2ce7g5nncgwidbsk2d67g5nmh5w9rimu2c87g5nmtck49bsk2d77g5nmvqohx9ac2cj7g5nme883frwi2cb7g5nnaigxvimu2cc7g5nn84cnrrwi2cd7g5nmjk4n59ac2cn7g5nnev43frwi2ch7g5nmvb4vghhc2cu7g5nmqugxg3e42c57g5nmiksrcraq2cv7g5nmc9oporaq2cw7g5nn93kx7pie2c97g5nmm28m3imu2cp7g5nnkqknwraq2d77g5nnefgeoraq2dd7g5nnla8r59ac2d27g5nn384b3imu2cw7g5nmh5wdoraq2dm7g5nnbxcrjrwi2db7g5nnaiga3rwi2dc7g5nmqespw49i2bd7g5nmm289rrwi2dh7g5nn3noep9ac2di7g5nn7kohghhc2do7g5nnqm4orrwi2dt7g5nncgwokraq2e97g5nmc9ornrwi2du7g5nn7korbimu2ds7g5nmqes7shhc2e77g5nnp78eg49i2cd7g5nmre4kjpie2dh7g5nnhssdcraq2eq7g5nnm9gfmt2e2e37g5nmm28vut2e2e67g5nnskko3imu2el7g5nnaigbshhc2f47g5nn3842brwi2et7g5nnp78kbpie2ed7g5nnn8olwraq2fp7g5nmhlkfvimu2f87g5nmdokaqt2e2er7g5nnn8o4wraq2fs7g5nml2w299ac2f97g5nn7kobs49i2dd7g5nmlik5mt2e2eu7g5nmfn4r8raq2g37g5nn84cgpbsk2gj7g5nn5m8aat2e2f67g5nngtk69bsk2gn7g5nnjrcucraq2gd7g5nn65ssjpie2f27g5nnqm47w3e42f97g5nmkjcn6t2e2fe7g5nnpmsbvpie2fc7g5nnefgo849i2e97g5nnr5obdbsk2h77g5nnm9g543e42fg7g5nmrtolchhc2h37g5nmm2868hhc2h47g5nn9j8st9ac2fx7g5nmmxgdbpie2ff7g5nmqesck49i2ea7g5nna2ssfrwi2g47g5nn9j8lchhc2hi7g5nmm28drimu2g47g5nnis4xfrwi2gb7g5nn6lg5s3e42fr7g5nnbhobg49i2es7g5nn9j8rdbsk2hw7g5nn56kgnrwi2gn7g5nnjrcu83e42gf7g5nn84c5qt2e2gm7g5nmk3oxdbsk2i97g5nmc9onfrwi2gu7g5nmvb4wpbsk2ic7g5nmubs2frwi2h37g5nmsdcbfimu2h27g5nna2sa43e42gt7g5nme88og3e42gu7g5nnlpsuw3e42h57g5nmgmcqdbsk2is7g5nmjk4d43e42he7g5nmliktet2e2h97g5nmj4gxk49i2fu7g5nmvqoisraq2ii7g5nnaigigraq2ik7g5nmvqo483e42hp7g5nmc9owqt2e2hk7g5nnla8plbsk2jk7g5nnfucfut2e2ho7g5nmtcknghhc2jq7g5nmx9km7imu2i57g5nmmhspo49i2gf7g5nno7w6rimu2i67g5nngdw9h9ac2il7g5nndvse5bsk2jv7g5nnefgm59ac2io7g5nmts8a7imu2ig7g5nmurg5kraq2j87g5nnfeokghhc2kd7g5nmxp83hbsk2kn7g5nnqm4v7imu2j37g5nngdwkxbsk2kv7g5nmubsjrrwi2ix7g5nmnwo243e42is7g5nnpmsbg3e42j57g5nmf7gxjpie2j57g5nn3nouc49i2i37g5nn5m8sfimu2jl7g5nmm28w9bsk2ln7g5nmmxgarpie2jd7g5nnev4ssraq2kp7g5nmkjce9bsk2lw7g5nn56k78raq2l57g5nnbxccat2e2ji7g5nnkawmc49i2ir7g5nnis4jk49i2iw7g5nnp78vvrwi2k77g5nmqesslbsk2m37g5nmre4fwraq2l87g5nmmhssw49i2j27g5nn75468raq2lc7g5nmhlkorimu2kb7g5nmx9kkpbsk2m77g5nn65shqt2e2js7g5nncgw7hbsk2md7g5nn93kwkhhc2mb7g5nncwklpbsk2mf7g5nmg6orwraq2ll7g5nmg6osshhc2mk7g5nmpv8knimu2kw7g5nn65sqet2e2kc7g5nngdwaoraq2lr7g5nnla8d6t2e2kd7g5nngtktrpie2kf7g5nngdw6nrwi2kn7g5nmf7grnrwi2kq7g5nmubs8bpie2kj7g5nnm9gibrwi2kx7g5nncwk92t2e2kh7g5nmj4gdtbsk2mw7g5nn8jwjd9ac2l77g5nn93kjbimu2ln7g5nn2ogtk3e42kv7g5nmnh4299ac2lk7g5nmnwo9jimu2lq7g5nnmp45g3e42l27g5nn84co849i2kc7g5nngdw53pie2l57g5nmvb4ret2e2kv7g5nmre4kk3e42la7g5nn2og5ghhc2nc7g5nnkqku7rwi2ls7g5nmpv8ret2e2l67g5nmk3odo3e42ll7g5nmovw8frwi2lu7g5nnhss9d9ac2m47g5nmi58pghhc2o47g5nmts88x9ac2m97g5nn56k3dbsk2od7g5nmensiit2e2lo7g5nnaigivpie2m37g5nnnoc8graq2nn7g5nns4whfimu2mv7g5nnmp4h43e42lw7g5nmquggdbsk2ow7g5nmgmcwpbsk2p27g5nmpv8pshhc2oq7g5nn93kpw49i2lw7g5nna2s78hhc2ot7g5nmd8w8craq2o87g5nml2wmx9ac2mq7g5nn28sm4raq2oc7g5nnev45vrwi2mx7g5nn65scchhc2p97g5nmvqots3e42n27g5nmm282vrwi2nh7g5nnjrcvc3e42n47g5nmenslut2e2mx7g5nnhd8wnimu2nv7g5nmgmctg3e42na7g5nncgwvdbsk2px7g5nn7ko8rpie2nc7g5nnqm4f99ac2nf7g5nnjbono49i2mw7g5nmsswjhbsk2qc7g5nnhssxhbsk2qf7g5nncwkpx9ac2nk7g5nnicg9p9ac2nl7g5nngtk7graq2ph7g5nnpmshet2e2nl7g5nnjbobet2e2nt7g5nmg6onit2e2nu7g5nmre46fimu2oe7g5nnbxc3rimu2og7g5nn8jwschhc2qf7g5nnla8lcraq2q57g5nmwpwcpbsk2r47g5nmnh4hl9ac2oi7g5nnaigwvimu2p27g5nnfuc4849i2np7g5nnfeoq3pie2om7g5nn65snrpie2ot7g5nn8jwnit2e2op7g5nmpfka2t2e2oq7g5nnr5o28hhc2rc7g5nn3nohhbsk2rr7g5nmiks3fpie2p77g5nmre4n4hhc2ri7g5nngdwmg3e42ou7g5nmf7g9graq2r57g5nmvqowjimu2pq7g5nmre4qfimu2px7g5nmrtovvpie2ph7g5nnhss7kraq2rj7g5nmd8wjet2e2pc7g5nn6lgcgraq2rq7g5nnkqkt9bsk2sk7g5nmqugkh9ac2pm7g5nnm9g3frwi2qv7g5nnmp4qrrwi2qx7g5nnaigvrpie2q77g5nme88nrrwi2r27g5nmmhsfk3e42pk7g5nnpmsbx9ac2pr7g5nns4wvat2e2q67g5nn7koknimu2qq7g5nnbxct3imu2qt7g5nmmxgf3pie2qe7g5nn28spjimu2qx7g5nnq6ghvrwi2ro7g5nmf7gl8hhc2tb7g5nnskk8d9ac2q57g5nmnwoh83e42qi7g5nmm28g4hhc2tf7g5nmxp87frwi2s37g5nn93kdfimu2rr7g5nmfn45k3e42qt7g5nmnwoerimu2rv7g5nn84c3d9ac2qn7g5nn65sl449i2ps7g5nngdwtrimu2rw7g5nngtkkg3e42r27g5nmvqof7rwi2sh7g5nnnoc3849i2qg7g5nmj4grnimu2sp7g5nn28sqw3e42rr7g5nmvb4n8hhc2ug7g5nnjrcbmt2e2s37g5nmre47c3e42s57g5nn56kbc3e42s87g5nnlpsgl9ac2ri7g5nmqes7ohhc2up7g5nmdokmit2e2sa7g5nnqm488hhc2uq7g5nnp786whhc2us7g5nnfeofrrwi2tb7g5nngtkccraq2uf7g5nnhsshw3e42ss7g5nnaigl849i2rd7g5nnfucfbrwi2tk7g5nmkjcdrpie2sm7g5nn93kpwraq2ur7g5nmfn45fimu2tk7g5nnjbokghhc2v87g5nmens8it2e2sn7g5nmf7gc8raq2us7g5nngdw83pie2sp7g5nmxp8jw49i2rj7g5nmnh4cxbsk2vb7g5nnq6g8849i2rk7g5nmmxgh9bsk2vc7g5nmsdcu83e42t97g5nngdw68raq2v27g5nmh5wbrrwi2u37g5nmqug7849i2s37g5nnqm4lvpie2t67g5nmj4gxet2e2ta7g5nmnh49ghhc2vj7g5nn4mwqc3e42th7g5nmgmc7ut2e2tj7g5nnhss8l9ac2t47g5nn28svdbsk2w47g5nnev4q4raq2vg7g5nmg6oso3e42tp7g5nns4wot9ac2t87g5nmsswxk3e42tw7g5nnis4gut2e2tt7g5nn6lghwraq2vv7g5nmmhs8oraq2w27g5nnmp4rkraq2w37g5nmovw5wraq2w67g5nmh5wcwraq2w77g5nmmhs4vpie2ua7g5nncwkn6t2e2ue7g5nmcpcmat2e2ug7g5nn47c3nimu2un7g5nnn8ofd9ac2tu7g5nmm28to49i2ti7g5nngdwqo3e42um7g5nmcpc8rrwi2vf7g5nmvqo7c49i2u27g5nnfuchfrwi2vj7g5nnr5ohut2e2vf7g5nncgwpnpie2vd7g5nngtkhoraq2xc7g5nmh5wek49i2uc7g5nnfeopkraq2xj7g5nmiksajpie2vl7g5nmj4g34raq2xm7g5nn93kwo49i2un7g5nmogci2t2e2vx7g5nn2ogkw3e42ve7g5nmmhsw83e42vf7g5nmxp8nt9ac2ut7g5nnb24o849i2v47g5nmts8v59ac2uw7g5nnmp4vbrwi2w97g5nnn8offimu2vl7g5nmnh4a4hhc3247g5nmogcwrrwi2wj7g5nmovw4k3e42vm7g5nmlikeohhc3287g5nmc9osnimu2w37g5nmj4gefrwi2wq7g5nmqugt8raq32k7g5nnhd8a43e42vv7g5nmc9oqc49i2vk7g5nn56kah9ac2vj7g5nmsswgfimu2we7g5nmlik7bpie2wv7g5nncgw8vrwi2x67g5nmikstc3e42w67g5nn28s9frwi2x87g5nmwacvmt2e2xb7g5nmwpwnbimu2wk7g5nmvb4eghhc32x7g5nnq6gd9bsk33o7g5nmmhs63imu2wx7g5nngdw9l9ac2w47g5nmcpcl83e42wf7g5nn4mw5nimu2x27g5nnlpsdbrwi2xn7g5nnn8oqjpie2xi7g5nnbhomjimu2x97g5nn9j8hc3e42wl7g5nmk3o49bsk3477g5nmlikw4raq33m7g5nnkawknrwi3277g5nmk3ohs49i2wu7g5nmiksfvpie32j7g5nmk3oik3e42xq7g5nmmxgdp9ac2x87g5nmcpcxbpie32u7g5nn754vh9ac2xk7g5nnhssdchhc34h7g5nnb244tbsk35k7g5nmkjcsvpie3357g5nn8jwinrwi33f7g5nmrtojl9ac2xq7g5nnis4ubimu3337g5nngtkjx9ac2xv7g5nmsdchx9ac3247g5nmvqorkhhc3587g5nmpfk4lbsk36e7g5nmsswrjimu33a7g5nn5m8srrwi33t7g5nmpv8uoraq3577g5nmubs93rwi3437g5nn384p7imu33k7g5nndvseo3e432k7g5nmi58lkraq35e7g5nn2ogsohhc35t7g5nmgmchl9ac32n7g5nme887c3e432t7g5nn7koup9ac32t7g5nmj4gpcraq35j7g5nnnoc7x9ac32w7g5nmnh4tdbsk3787g5nmwpwfbimu3447g5nnbhou2t2e34l7g5nmvb4thbsk37d7g5nmqugpjimu3467g5nmliko4hhc3687g5nn65s6npie34d7g5nn9j859bsk37l7g5nnq6gojimu34f7g5nnlpstkhhc36f7g5nmmhs9sraq3637g5nmiks2dbsk37s7g5nmcpcl8raq3667g5nmpv8wjrwi35a7g5nn84cvvrwi35b7g5nnla8xoraq36a7g5nnaigmvrwi35i7g5nnkqkmjrwi35m7g5nmmhst6t2e35j7g5nn4mwnw49i33p7g5nmjk4gnimu34v7g5nmvqoud9ac3457g5nnp783qt2e35u7g5nmwpw7whhc3767g5nnonkdbpie3587g5nns4w6bpie35e7g5nmcpccjpie35g7g5nmfn4ovrwi36h7g5nmi58id9ac34j7g5nnqm43rpie35m7g5nmk3o63imu35k7g5nn5m8eo3e434i7g5nnkqk3x9ac34p7g5nn754qnpie35s7g5nndvsr83e434u7g5nmvqotwhhc3887g5nnq6gsxbsk39a7g5nns4wosraq37s7g5nno7w53rwi3777g5nmkjc93imu36f7g5nnkawd7imu36j7g5nmj4g5jimu36k7g5nnefgi4raq3867g5nnkawdlbsk39p7g5nndg8ffpie36n7g5nnhd8xat2e3787g5nmovwvfpie36u7g5nnr5o3rimu3757g5nmogc3whhc38w7g5nndg8knimu37b7g5nnaigl59ac36i7g5nmx9k243e436c7g5nn56kivimu37d7g5nmmhsxwhhc3967g5nmiks8jpie37i7g5nmgmcqdbsk3aq7g5nnpmsmo49i36o7g5nnefgnshhc39k7g5nmubsow49i36p7g5nnm9g9rpie37p7g5nmwacd849i36q7g5nno7wb3imu37p7g5nns4wig49i36r7g5nnlpsgoraq38v7g5nmcpc42t2e38d7g5nn8jwj3rwi38l7g5nml2wmh9ac3797g5nnp7836t2e38j7g5nmtckajimu37v7g5nmjk46w3e436p7g5nnjrcgkraq3977g5nn3no87pie38c7g5nn7koqsraq39f7g5nmkjc7xbsk3bk7g5nngdwajimu38n7g5nnskkbs3e437d7g5nnfuck3imu38p7g5nn56kol9ac3827g5nna2sb7rwi39h7g5nn9j8ebpie3987g5nnp784it2e39b7g5nmnwoipbsk3c57g5nmj4g4fimu38u7g5nml2ww59ac3877g5nnb24983e437u7g5nmtckarpie39m7g5nn3noh7imu3937g5nn4mw6fimu3967g5nmvqosmt2e39l7g5nn5m84w49i38k7g5nno7wjshhc3bt7g5nncgwpfpie3a57g5nnskkuk49i38l7g5nn9j85o49i38n7g5nnq6guvimu39k7g5nmh5w8mt2e39u7g5nn56kvd9ac3947g5nncgwrh9ac3957g5nmcpc3bimu39r7g5nmmxgjkhhc3c57g5nn9j8ppbsk3cu7g5nmsswel9ac39a7g5nmiksinpie3ah7g5nnaigpo49i3927g5nns4wrvimu3a87g5nmkjccdbsk3d87g5nnpms8rrwi3ax7g5nmvqolx9ac39j7g5nn65sdwraq3b87g5nml2wnrpie3bj7g5nmcpcvohhc3de7g5nmf7gjo3e43a27g5nnn8on9bsk3dv7g5nnev4dxbsk3dw7g5nnhssa9bsk3e97g5nmwpw37imu3bc7g5nmkjc4g3e43am7g5nmfn4tghhc3e67g5nnla8lkhhc3eb7g5nn4mwlk3e43as7g5nmfn4o449i3ar7g5nmmxg2xbsk3el7g5nn28s4nimu3bx7g5nmtcke99ac3ba7g5nn3noigraq3cl7g5nmqugiw49i3b67g5nme885rpie3cn7g5nnaighet2e3co7g5nn5m8ol9ac3bq7g5nngtkskhhc3fi7g5nmj4g2ut2e3cv7g5nmfn4o43e43bk7g5nn65sq59ac3bw7g5nmcpcd7rwi3dx7g5nnfucls49i3c27g5nmcpciit2e3d97g5nnqm4pwraq3db7g5nn4mw2fimu3cq7g5nn28sts3e43bu7g5nmpfkss49i3cn7g5nnskkhw49i3cp7g5nmurgfohhc3gd7g5nnskkjsraq3ds7g5nmjk4m4hhc3gj7g5nngtkkrimu3dg7g5nmx9k7whhc3gm7g5nmvqofkhhc3gp7g5nn9j88craq3e37g5nmvb4rcraq3e77g5nnpmsbcraq3e87g5nn56kqit2e3dx7g5nmi58ec49i3dj7g5nnicg4ohhc3gw7g5nnjrcarimu3dv7g5nmf7ghg3e43cu7g5nnis4mlbsk3gp7g5nmjk4lqt2e3eh7g5nnlps5shhc3ht7g5nmhlk3craq3f47g5nn3nov7imu3eg7g5nmj4gk43e43db7g5nnpms6wraq3f87g5nmts8rit2e3f57g5nn754sl9ac3dl7g5nn9j8bjpie3fu7g5nmts8r3imu3er7g5nmf7gn3imu3ew7g5nnkqk7shhc3ig7g5nmpv8fw3e43du7g5nmcpctqt2e3fr7g5nmi587mt2e3fs7g5nnis4ctbsk3ic7g5nnb24kpbsk3ii7g5nmx9kjbpie3gw7g5nnmp4u59ac3eb7g5nmrtoxrimu3fr7g5nmnh45tbsk3ip7g5nmnh4tsraq3g87g5nmogcsoraq3ga7g5nno7whrimu3g97g5nnhsscmt2e3gp7g5nmcpc5pbsk3j27g5nmcpciohhc3jk7g5nn384wo49i3g97g5nmg6oo8raq3gk7g5nmfn4s849i3ga7g5nn93kuohhc3jm7g5nmj4ge43e43eu7g5nnr5o4brwi3ha7g5nn5m8dwhhc3jn7g5nnbhofo3e43ff7g5nnn8or7rwi3hk7g5nmubsgfpie3ia7g5nmc9opx9ac3ff7g5nnp788l9ac3fi7g5nmurgesraq3hn7g5nnqm4b2t2e3hj7g5nnrlc4vpie3ip7g5nmxp8vxbsk3jq7g5nnicgx83e43g37g5nmg6ob4hhc3kj7g5nmpfkf7rwi3ib7g5nnb24tc49i3hk7g5nn6lgl849i3hm7g5nn5m8cmt2e3i97g5nmhlk78raq3ic7g5nmcpc6jrwi3ip7g5nmurgx83e43gs7g5nnfeo2x9ac3gn7g5nnn8o5s49i3hu7g5nnjbo5k49i3hw7g5nmwac6bpie3jm7g5nnev4443e43gx7g5nmwact3rwi3ix7g5nmj4g8s3e43h57g5nmpfkm3rwi3j37g5nmts8bvimu3hs7g5nmurg6craq3iq7g5nnla86chhc3ll7g5nmxp8ek3e43hb7g5nmubs88raq3iu7g5nns4w943e43hg7g5nmnh4n7pie3k27g5nnq6g9d9ac3hj7g5nnqm4jfpie3k77g5nnaigiit2e3jc7g5nmtckdwraq3ja7g5nnis4set2e3ji7g5nmfn4h3rwi3jm7g5nmpv8pbpie3kj7g5nnkawqw3e43hv7g5nnqm4ejrwi3jr7g5nn7kolfrwi3jt7g5nnicgkrimu3in7g5nnnocerrwi3k47g5nnfeoap9ac3iq7g5nmwacg6t2e3k47g5nmjk4qc49i3jd7g5nmd8wl6t2e3k67g5nnbhooet2e3k87g5nmm28u99ac3j47g5nnnocdl9ac3j77g5nmdokanimu3j47g5nme88msraq3k77g5nmc9op8raq3kd7g5nn9j8kvimu3jb7g5nnis47qt2e3kn7g5nnmp49vpie3lg7g5nmjk4lxbsk3mr7g5nnlpsd2t2e3l27g5nns4wop9ac3ju7g5nnb24e3pie3ln7g5nmogco7rwi3l57g5nn65sd99ac3k87g5nn56kro3e43jk7g5nnjbon849i3kl7g5nmrto2849i3km7g5nnaig2k3e43jo7g5nmrtovrrwi3lm7g5nn8jw6k3e43jv7g5nnbhoo8raq3lf7g5nmurgg59ac3km7g5nn47cncraq3lj7g5nnn8okrrwi3ma7g5nmgmcovpie3mu7g5nmk3ols49i3le7g5nn7koqs3e43k97g5nnmp4j7rwi3ml7g5nnfeolfpie3mx7g5nmg6o7w49i3lk7g5nmxp8jbimu3kr7g5nmsdchut2e3me7g5nn2ogmjimu3kt7g5nmx9kanpie3na7g5nmcpctfrwi3mw7g5nnfuc399ac3ld7g5nmxp8qw49i3m47g5nmensaw3e43kq7g5nmc9o6nimu3l97g5nmf7gn5bsk3p77g5nn7545h9ac3lx7g5nnb24brimu3li7g5nnicgvvimu3ll7g5nnonkhs3e43la7g5nna2s259ac3mc7g5nmvqo3dbsk3pg7g5nmi582pbsk3ph7g5nnicglchhc3pu7g5nmxp8os49i3n47g5nmts83bimu3m27g5nmi58g4raq3nl7g5nnkqkgjpie3oo7g5nmhlkw8hhc3q67g5nmtckqut2e3nq7g5nnev4q3rwi3o97g5nnskkd7rwi3oa7g5nn8jws43e43m37g5nme88v8raq3nq7g5nmcpc383e43me7g5nnkqkg3pie3p77g5nmk3oc9bsk3q87g5nnjrcfd9ac3nb7g5nnr5ol8raq3o97g5nn7ko64hhc3qn7g5nnnoc383e43mt7g5nmk3opjrwi3p87g5nmkjchat2e3od7g5nngdwo8hhc3qs7g5nnhss5mt2e3oi7g5nmdokr7rwi3pe7g5nncgwg9bsk3qu7g5nna2st7rwi3pj7g5nnla8jbimu3nd7g5nmnh4d8hhc3rc7g5nmqessvimu3ng7g5nmhlkwkraq3ow7g5nnkqkug3e43ng7g5nme88fk49i3ox7g5nn4mwlwraq3p47g5nnaigi3imu3nj7g5nnonklsraq3p77g5nn384n99ac3o97g5nmh5wemt2e3p37g5nmubssrpie3qd7g5nnbhoos3e43o87g5nnfuckh9ac3oq7g5nn84c6x9ac3os7g5nmdokix9ac3ot7g5nnjrct849i3pm7g5nmurgbw49i3pp7g5nmkjcnwhhc3sx7g5nnbho8sraq3q97g5nnhd8xlbsk3si7g5nncgw6nimu3op7g5nmtck4bpie3rb7g5nmensnbimu3oq7g5nnlpsit9ac3pm7g5nmensvbrwi3rc7g5nnis4wxbsk3sp7g5nmwac93rwi3rf7g5nnb24t59ac3ps7g5nme88j3imu3ox7g5nmlika8hhc3ts7g5nn384cohhc3tv7g5nnkqkawraq3qn7g5nmts8mo3e43ps7g5nmh5wqd9ac3q77g5nmdokb849i3r67g5nmqugvkraq3qu7g5nn384ns3e43px7g5nmfn4gtbsk3tc7g5nnfeo7c49i3rf7g5nme88j83e43q37g5nns4wkc49i3rh7g5nnonkxfpie3s27g5nmogcdit2e3qp7g5nmubs24raq3r97g5nnfeo63rwi3se7g5nmx9k7qt2e3qv7g5nmqes54raq3ri7g5nme88jkraq3rj7g5nmj4gt83e43qc7g5nnnoc8449i3rt7g5nml2w2qt2e3r47g5nnp78tw49i3rw7g5nn8jwhjpie3s67g5nn4mws8raq3rr7g5nnjboiohhc3uo7g5nnicgqet2e3rf7g5nmkjckat2e3rj7g5nmogcdl9ac3qx7g5nmiksi8raq3s77g5nndvs8fimu3q47g5nmtckarpie3sh7g5nmh5w9fimu3qb7g5nn5m8x43e43r27g5nnev4ulbsk3ub7g5nnfeo2d9ac3re7g5nmgmcol9ac3ri7g5nnhd8ggraq3so7g5nnmp4a3rwi3tg7g5nnn8op3pie3ss7g5nn2ogmdbsk3up7g5nml2wih9ac3rp7g5nmsdcfcraq3sv7g5nml2wwjrwi3to7g5nnkqkhlbsk3uw7g5nncgw9mt2e3si7g5nnpmsbqt2e3sk7g5nngdwx6t2e3sm7g5nmmhsek3e43s47g5nnfucmnrwi3u77g5nn3noi2t2e3sw7g5nn2ogkwhhc3wd7g5nme889shhc3wj7g5nn6lglbpie3u57g5nnb24uh9ac3sm7g5nn9j8a849i3u97g5nnonke7pie3ud7g5nn93kpjrwi3v77g5nmrto4pbsk3wb7g5nnrlckbrwi3vi7g5nnskkr3rwi3vn7g5nmliksxbsk3wm7g5nnskkggraq3um7g5nmogccchhc3xn7g5nnkqk55bsk3wv7g5nnicg5fimu3sp7g5nmqugxg3e43u67g5nn84cusraq3uv7g5nmkjcuxbsk3x27g5nmd8wgqt2e3uf7g5nmmxgrc49i3vh7g5nmre4rh9ac3u97g5nnfeogfrwi3wp7g5nmpv8nchhc42n7g5nnkawqshhc42o7g5nmensxbrwi3ws7g5nml2wat9ac3uf7g5nnlpssshhc42t7g5nnpms5vimu3tm7g5nn3nofqt2e3uw7g5nmogcwghhc42w7g5nmk3odbrwi3x67g5nnev4iwhhc4367g5nmjk4m7imu3u27g5nnhd8m7rwi3xb7g5nnjrcng49i3w87g5nnq6g6vimu3u97g5nmwacabimu3uj7g5nmgmc3graq3w97g5nnbhoecraq3wa7g5nn8jwwbimu3ur7g5nmwpwkat2e3vu7g5nno7w43pie3xj7g5nmubs4x9ac3vi7g5nmmhsapbsk42i7g5nmxp8n7pie3xn7g5nmhlkxw3e43w37g5nmiks35bsk42m7g5nmnwoapbsk42o7g5nmhlknnpie3xr7g5nmsdcjp9ac3vw7g5nmwpwoh9ac3w27g5nnkaw58raq3wn7g5nml2wgdbsk42s7g5nmovwtk49i3xf7g5nnefgxjimu3vb7g5nnefgbh9ac3w97g5nnaigvfpie423`

func minify() error {
	g, err := makeGenerator()
	if err != nil {
		return err
	}

	id := g.New(123)
	text, err := id.MarshalText()
	if err != nil {
		return err
	}
	data := string(text)

	fmt.Println("input:\n", data)

	var b bytes.Buffer

	// Compress the data using the specially crafted dictionary.
	zw, err := flate.NewWriterDict(&b, flate.DefaultCompression, []byte(dict))
	if err != nil {
		return err
	}
	if _, err := io.Copy(zw, strings.NewReader(data)); err != nil {
		return err
	}
	if err := zw.Close(); err != nil {
		return err
	}

	fmt.Printf("length: %d\n", len(data))
	fmt.Printf("length: %d\n", len(b.String()))

	// The decompressor must use the same dictionary as the compressor.
	// Otherwise, the input may appear as corrupted.
	fmt.Println("Decompressed output using the dictionary:")
	zr := flate.NewReaderDict(bytes.NewReader(b.Bytes()), []byte(dict))
	if _, err := io.Copy(os.Stdout, zr); err != nil {
		return err
	}
	if err := zr.Close(); err != nil {
		return err
	}

	fmt.Println()

	// Substitute all of the bytes in the dictionary with a '#' to visually
	// demonstrate the approximate effectiveness of using a preset dictionary.
	fmt.Println("Substrings matched by the dictionary are marked with #:")
	hashDict := []byte(dict)
	for i := range hashDict {
		hashDict[i] = '#'
	}
	zr = flate.NewReaderDict(&b, hashDict)
	if _, err := io.Copy(os.Stdout, zr); err != nil {
		return err
	}
	if err := zr.Close(); err != nil {
		return err
	}

	return nil
}
