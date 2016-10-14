import { showModal } from './modal.js'

function buildSignupHtml() {
    var html = `
        <h1>Create an account</h1>
        <label for="username">Username</label>
			     <input type="text" class="full" name="username"/>
			     <label for="password">Password</label>
			     <input type="password" class="full" name="password"/>
        <input type="button" value="Sign Up" class="signup-button"/>
    `
    return html;
}

function onSignup(callback) {
    var $usernameField = $('input[name=username]')
    var $passwordField = $('input[name=password]')
    var username = $usernameField.val()
    var password = $passwordField.val()
    // TODO: validation
    $.post('/signup', {
        username: username,
        password: password
    }, callback)
           
}

export function showLogin() {

}

export function showSignup() {
    var modal = showModal({
        content: buildSignupHtml()
    });
    $('input.signup-button').click(onSignup)
}
