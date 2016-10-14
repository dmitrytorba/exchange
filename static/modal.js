
function buildHtml(content) {
    var html = `
        <div class="modal">
           <div class="modal-body">
           ${content}
           </div>
        </div>
        <div class="modal-overlay">
        `

    return html
}

function closeModal($modal) {
    $modal.hide()
}

export function showModal(config) {
    var content = config.content || ''
    var $modal = $(buildHtml(content))    
    $('body').append($modal)
    $('.modal-overlay').click(e => closeModal($modal)) 
}
