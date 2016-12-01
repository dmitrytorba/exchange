import techan from 'techan'

var margin = {
    top: 20,
    right: 20,
    bottom: 30,
    left: 50
}

function getTrades() {
    var url = '/bitfinex/trades/btcusd'
    return $.get(url)
}

function getSvg(selector) {
    var $container = $(selector)
    if ($container.has('svg').length) {
        //TODO
    } else {
        var svg = d3.select(selector).append('svg')
        svg.attr('height', $container.height())
        svg.attr('width', $container.width())
        svg.width = $container.width() - margin.left - margin.right
        svg.height = $container.height() - margin.top - margin.bottom

        svg.g = svg.append('g')
        svg.g.attr("transform", "translate(" + margin.left + "," + margin.top + ")")

        return svg;
    }
}

// line chart: https://bl.ocks.org/mbostock/3883245
function drawLineChart(data, selector) {
    var svg = getSvg(selector)
    var g = svg.g
    var width = svg.width
    var height = svg.height

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
