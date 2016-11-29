
function getTrades() {
    var url = '/bitfinex/trades/btcusd'
    return $.get(url)
}

function drawCandlestick(data, selector) {
    //TODO
    
}

export function drawPriceChart(selector) {
    getTrades().done((data) => {
        drawCandlestick(JSON.parse(data), selector)
    })
}
