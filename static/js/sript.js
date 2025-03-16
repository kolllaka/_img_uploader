const toast = new Toast()

const getItemTemplate = (link, id) => {
	return `
		<div class="item">
			<p class="item__filename">${link}</p>
			<img src="${link}" class="item__img">
			<div class="item__btns">
				<div data-btn="copy" class="item__btn btn">copy</div>
				<div data-btn="del" data-id="${id}" class="item__btn btn">del</div>
			</div>
		</div>
	`
}

const gallery = document.querySelector('.gallery')

gallery.innerHTML = data.map((img) => {
	return getItemTemplate(img.path, img.id)
}).reverse().join('')


document.addEventListener('click', (e) => {
	if (e.target.closest('.btn')) {
		switch (e.target.closest('.btn').dataset.btn) {
			case "copy":
				copyHandler(e)

				break;
			case "del":
				deleteHandler(e)

				break;

			case "link":
				linkHandler(e)

				break;

			case "file":
				fileHandler(e)

				break;
			default:
				console.log('unknow btn!');

				break;
		}
	}
})

const copyHandler = (e) => {
	const filepath = e.target.closest('.item').querySelector('.item__filename').textContent
	const url_link = document.getElementById('url_link').value

	if (url_link === '') {
		toast.create('pls fill "url link" field', 'error')

		return
	}

	const copyText = url_link + filepath

	navigator.clipboard.writeText(copyText)
		.then(() => {
			toast.create("link copy", "ok")
		})
		.catch(err => {
			toast.create("sorry your browser dont allow copy command", "error")

			console.error(err)
		})
}

const deleteHandler = (e) => {
	console.log('delete');
	console.log(e.target)

	const id = e.target.dataset.id

	fetch(`api/delete?id=${id}`, {
		method: "DELETE",
	})
		.then((response) => {
			switch (response.status) {
				case 204:
					toast.create('image delete', 'ok')

					data = data.reverse().filter(img => img.id !== id)

					gallery.innerHTML = data.map((img) => {
						return getItemTemplate(img.path, img.id)
					}).reverse().join('')

					return
				case 400:
					throw new Error('Bad request')
				default:
					throw new Error('Something went wrong')
			}
		})
		.catch((err) => {
			toast.create('bad request', 'error')

			console.error(err)
		});
}

const linkHandler = (e) => {
	const input = e.target.parentElement.querySelector('input')
	const link = input.value

	if (link === '') {
		toast.create('pls fill link field', 'error')

		return
	}


	body = JSON.stringify({ link: link })

	fetch('api/link', {
		method: "POST",
		headers: {
			"Content-Type": "application/json",
		},
		body: body,
	})
		.then((response) => {
			switch (response.status) {
				case 200:
					return response.json()
				case 400:
					throw new Error('Bad request')
				default:
					throw new Error('Something went wrong')
			}
		})
		.then((resp) => {
			console.log(resp);
			data.push(resp)

			gallery.insertAdjacentHTML('afterbegin', getItemTemplate(resp.path, resp.id))
		})
		.catch((err) => {
			toast.create('bad request', 'error')

			console.error(err)
		})

	input.value = ''

	console.log('upload')
}

const fileHandler = (e) => {
	const input = e.target.parentElement.querySelector('input')

	if (input.files.length < 1) {
		toast.create("pls choose file to upload", "error")

		return
	}

	const image = input.files[0]
	const formData = new FormData()
	formData.append("image", image)

	fetch('api/upload', { method: "POST", body: formData })
		.then((response) => {
			switch (response.status) {
				case 200:
					return response.json()
				case 400:
					throw new Error('Bad request')
				default:
					throw new Error('Something went wrong')
			}
		})
		.then((resp) => {
			console.log(resp);
			data.push(resp)

			gallery.insertAdjacentHTML('afterbegin', getItemTemplate(resp.path, resp.id))
		})
		.catch((err) => {
			toast.create("Bad request", "error")
			console.error(err)
		})

	input.value = ''
	e.target.parentNode.querySelector('.file__placeholder').style.display = "block"
	e.target.parentNode.querySelector('.file__text').style.display = "none"
	e.target.parentNode.classList.remove('choosen')
}

// files input
const fileInputs = document.querySelectorAll('input[type=file]')

Array.prototype.forEach.call(fileInputs, (input) => {
	input.addEventListener('change', (e) => {
		const length = e.target.files.length

		if (length > 0) {
			toast.create(`file selected`, "ok")
			e.target.parentNode.querySelector('.file__placeholder').style.display = "none"
			e.target.parentNode.querySelector('.file__text').style.display = "block"
			e.target.parentNode.classList.add('choosen')
		} else {
			e.target.parentNode.querySelector('.file__placeholder').style.display = "block"
			e.target.parentNode.querySelector('.file__text').style.display = "none"
			e.target.parentNode.classList.remove('choosen')
		}
	})
})
