import { showModal } from './modal.js'

function buildLoginHtml() {
    var html = `
        <h1>Welcome</h1>
        <div class="error-feedback"></div>
        <label for="username">Username</label>
			     <input type="text" class="full username-field" name="username"/>
			     <label for="password">Password</label>
			     <input type="password" class="full password-field" name="password"/>
        <input type="button" value="Log In" class="login-button"/>
    `
    return html;
}

function buildNav(username) {
    return `
        <a href="/settings" class="account-button">
            ${username}
        </a>   
    `
}

function onLogin(modal) {
    var $usernameField = $('.username-field', modal.$el)
    var $passwordField = $('.password-field', modal.$el)
    var username = $usernameField.val()
    var password = $passwordField.val()
    // TODO: validation
    // TODO: csrf
    $.post('/login', {
        username: username,
        password: password
    })
    .done((user) => {
        $('.login-button').hide()
        $('.signup-button').hide()
        $('.header .nav').html(buildNav(user))
        modal.closeModal()
    })
    .fail(() => {
        $('.error-feedback', modal.$el).text('Incorrect login.')
    })
}

export function showLogin() {
    var modal = showModal({
        content: buildLoginHtml()
    });
    $('input.login-button', modal.$el).click(event =>
                                             onLogin(modal))
}
