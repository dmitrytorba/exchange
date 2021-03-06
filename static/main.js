import $ from 'jquery'
import 'flot'
import { login, logout } from './login.js'
import { signup } from './signup.js'
import { drawPriceChart } from './charts.js'
import { drawStats } from './stats.js'

//drawBook('#books', 'bitfinex')

drawStats('.gdax-stats-btcusd', 'gdax', 'btcusd')
drawStats('.gdax-stats-ethusd', 'gdax', 'ethusd')
drawStats('.gdax-stats-ethbtc', 'gdax', 'ethbtc')

// TODO: add router instead of click events
// $('body').on('click', '.header .login-button', login);
// $('body').on('click', '.header .logout-button', logout);
// $('body').on('click', '.header .signup-button', signup);

$(".tabs span").click(function(event){
	var parent = $(this).parent();
	var panel = parent.parent();
	var newtab = $("." + event.target.id + "-tab");

	// make sure the active tab is set
	panel.find(".active").removeClass("active");
	$(this).addClass("active");

	// make sure the active tab page is shown
	  newtab.addClass("active");
    if (event.target.id === "price") {
        drawPriceChart(".price-tab")
    }
});
