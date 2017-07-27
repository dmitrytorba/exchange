

function buildStatsHtml(title) {
    var html = `
        <div class="container">
           <h4>${title}</h4>
        </div>
        `

    return html
}

export function drawStats(selector, exchange, currency) {
    var $container = $(selector)
    $container.append(buildStatsHtml(''))
    var price = d3.select(selector).select('h4')
    var evtSource = new EventSource("/stats/" + exchange + "/" + currency)
    evtSource.onmessage = function(e) {
        var stats = JSON.parse(e.data)
        price.text(stats.Price)
    }
}
