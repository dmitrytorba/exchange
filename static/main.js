import $ from 'jquery'
import 'flot'
import {showLogin, showSignup} from './login.js'

// TODO: add router instead of click events
$('.login-button').click(showLogin);
$('.signup-button').click(showSignup);

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
	{ label: "Buys", data: buy_graph, color: "green" },
	{ label: "Sells", data: sell_graph, color: "red" },
], options);

$(".tabs span").click(function(event){
	var parent = $(this).parent();
	var panel = parent.parent();
	var newtab = $("." + event.target.id + "-tab");

	// make sure the active tab is set
	panel.find(".active").removeClass("active");
	$(this).addClass("active");

	// make sure the active tab page is shown
	newtab.addClass("active");
});
