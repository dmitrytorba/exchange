import { showModal } from './modal.js'

function buildSignupHtml() {
    var html = `
        <form>
        <h1>Create an account</h1>
        <div class="error-feedback"></div>
        <label for="username">Username</label>
			     <input type="text" class="full username-field" name="username" autofocus="autofocus"/>
			     <label for="password">Password</label>
			     <input type="password" class="full password-field" name="password"/>
        <input type="submit" value="Sign Up" class="signup-button"/>
        </form>
    `
    return html;
}

function submitSignup(modal) {
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


export function signup() {
    var modal = showModal({
        content: buildSignupHtml()
    })
    $('form', modal.$el).submit(event =>
                                        submitSignup(modal))
}
