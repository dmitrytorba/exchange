

function buildBookHtml(title) {
    var html = `
        <div class="container">
           <h4>
           ${title}
           </h4>
        </div>
        `

    return html
}

export function drawBook(selector, exchange) {
    var $container = $(selector)
    $container.append(buildBookHtml(exchange))
    var price = d3.select(selector).select('h4')
    var evtSource = new EventSource("/books/" + exchange)
    evtSource.onmessage = function(e) {
        var trade = JSON.parse(e.data)
        price.text(trade.Price)
    }
}
