import techan from 'techan'

function getTrades() {
    var url = '/bitfinex/trades/btcusd'
    return $.get(url)
}

function drawLineChart(data, selector) {
    var svg = d3.select(selector)
    var margin = {
        top: 20,
        right: 20,
        bottom: 30,
        left: 50
    }
    var width = +svg.attr("width") - margin.left - margin.right
    var height = +svg.attr("height") - margin.top - margin.bottom
    var g = svg.append("g")
        .attr("transform", "translate(" + margin.left + "," + margin.top + ")")

    var x = d3.scaleTime()
        .rangeRound([0, width]);

    var y = d3.scaleLinear()
        .rangeRound([height, 0]);

    var line = d3.line()
        .x(function(d) { return x(new Date(d.timestamp)) })
        .y(function(d) { return y(d.price) })

    x.domain(d3.extent(data, function(d) { return new Date(d.timestamp) }))
    y.domain(d3.extent(data, function(d) { return d.price }))

    g.append("g")
        .attr("class", "axis axis--x")
        .attr("transform", "translate(0," + height + ")")
        .call(d3.axisBottom(x));

    g.append("g")
        .attr("class", "axis axis--y")
        .call(d3.axisLeft(y))
        .append("text")
        .attr("fill", "#000")
        .attr("transform", "rotate(-90)")
        .attr("y", 6)
        .attr("dy", "0.71em")
        .style("text-anchor", "end")
        .text("Price (USD)");

    g.append("path")
        .datum(data)
        .attr("class", "line")
        .attr("d", line);
}

export function drawPriceChart(selector) {
    getTrades().done((data) => {
        drawLineChart(JSON.parse(data), selector)
    })
}
