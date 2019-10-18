.PHONY: setup geth

geth: setup
	geth \
		--datadir="./data" \
		--testnet \
		--syncmode=light \
		--ws \
		--wsport=8546 \
		--wsaddr=localhost \
		--wsorigins=statusjs,http://localhost:3000 \
		--rpc \
		--maxpeers=25 \
		--shh \
		--wsapi=web3,shh,admin,

setup:
	mkdir -p ./data/geth
	echo -e "[\n\
	\"enode://19872f94b1e776da3a13e25afa71b47dfa99e658afd6427ea8d6e03c22a99f13590205a8826443e95a37eee1d815fc433af7a8ca9a8d0df7943d1f55684045b7@35.238.60.236:30305\", \n\
	\"enode://960777f01b7dcda7c58319e3aded317a127f686631b1702a7168ad408b8f8b7616272d805ddfab7a5a6bc4bd07ae92c03e23b4b8cd4bf858d0f74d563ec76c9f@47.52.58.213:30305\" \n\
]" > ./data/geth/static-nodes.json

