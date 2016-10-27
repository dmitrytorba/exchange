
function buildHtml(content) {
    var html = `
        <div class="modal">
           <div class="modal-close">×</div>
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

    // since we appended the modal the DOM is technically not ready for any more actions afterwards
    // this seems kinda ghetto, can this be fixed?
    $(document).ready(function(){
        $modal.find("input:first").focus();
    });

    $('.modal-overlay').click(e => closeModal($modal))
    $('.modal-close').click(e => closeModal($modal)) 
    return {
        $el: $modal,
        closeModal: () => { closeModal($modal) }
    }
}
