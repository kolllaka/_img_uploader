class Toast {
	constructor(timeout = 3000) {
		this.$el = this.#init()
		this.timeout = timeout
	}

	#init() {
		const toast = document.createElement('div')
		toast.classList.add('toast')

		document.body.appendChild(toast)

		return toast
	}

	create(text, result = "ok" || "error") {
		const infoBody = this.#getTemplateToast(text, result)
		this.$el.insertAdjacentElement('beforeend', infoBody)

		setTimeout(function () {
			infoBody.remove()
		}, this.timeout);
	}

	#getTemplateToast(text, result) {
		const infoBody = document.createElement("div")
		infoBody.classList.add('toast__body')
		infoBody.classList.add(`toast__body-${result}`)
		infoBody.innerHTML = `<span>${text}</span>`

		return infoBody
	}
}


