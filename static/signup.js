import { showModal } from './modal.js'

function buildSignupHtml() {
    var html = `
        <h1>Create an account</h1>
        <div class="error-feedback"></div>
        <label for="username">Username</label>
			     <input type="text" class="full username-field" name="username"/>
			     <label for="password">Password</label>
			     <input type="password" class="full password-field" name="password"/>
        <input type="button" value="Sign Up" class="signup-button"/>
    `
    return html;
}

function onSignup(modal) {
    var $usernameField = $('.username-field', modal.$el)
    var $passwordField = $('.password-field', modal.$el)
    var username = $usernameField.val()
    var password = $passwordField.val()
    // TODO: validation
    // TODO: csrf
    $.post('/signup', {
        username: username,
        password: password
    })
    .done(() => {
        modal.closeModal()
    })
    .fail(() => {
        $('.error-feedback', modal.$el).text('Signup failed')
    })       
}


export function showSignup() {
    var modal = showModal({
        content: buildSignupHtml()
    });
    $('input.signup-button', modal.$el).click(event =>
                                              onSignup(modal))
}
