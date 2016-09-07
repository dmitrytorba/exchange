var buy_graph = [];
var sell_graph = [];

var last;
for (var i = 0; i < buys.length; i++) {
	var buy = buys[i];

	if (i > 0) {
		buy.amount += last.amount
	}

	buy_graph.unshift([buy.price, buy.amount]);

	last = buy;
}

for (var i = 0; i < sells.length; i++) {
	var sell = sells[i];

	if (i > 0) {
		sell.amount += last.amount
	}

	sell_graph.push([sell.price, sell.amount]);

	last = sell;
}

var options = {
    series: {
        lines: { show: true, fill: true }
    }
};

$("#spread-graph").plot([
	{ label: "Buys", data: buy_graph },
	{ label: "Sells", data: sell_graph },
], options);