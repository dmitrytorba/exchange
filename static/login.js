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
    if (username) {
        return `
        <a href="/settings" class="account-button">
            ${username}
        </a>
        <a href="#logout" class="logout-button">
            logout
        </a>   
        `
    } else {
        return `
        <a href="/login" class="login-button">
            login
        </a>
        <a href="/signup" class="signup-button">
            signup
        </a>
        `
    }
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

export function login() {
    var modal = showModal({
        content: buildLoginHtml()
    });
    $('input.login-button', modal.$el).click(e =>
                                             onLogin(modal))
    $('input', modal.$el).keyup(e => {
        if (e.keyCode === 13) {
            onLogin(modal)
        }
    })
}

export function logout() {
    $.get('/logout')
        .done(() => $('.header .nav').html(buildNav()))
}
