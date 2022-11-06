document
    .getElementById('check-availability-button')
    .addEventListener('click', (e) => {
        e.preventDefault()
        let html = `
                    <form id="check-availability-form" action="" method="post" novalidate class="needs-validation">
                        <div class="row">
                            <div class="col">
                                <div class="row" id="reservation-dates-modal">
                                    <div class="col">
                                        <input required class="form-control" type="text" name="start" id="start" placeholder="Check-in date" autocomplete="off"/>
                                    </div>
                                    <div class="col">
                                        <input required class="form-control" type="text" name="end" id="end" placeholder="Check-out date" autocomplete="off" />
                                    </div>
                                </div>
                            </div>
                        </div>
                    </form>
                    `
        notification.showDatepicker({
            title: 'Choose your dates',
            msg: html,
            callback: async (result) => {
                const form = document.getElementById('check-availability-form')
                const formData = new FormData(form)
                formData.append('csrf_token', '{{.CSRFToken}}')

                const response = await axios.post(
                    '/search-availability-json',
                    formData
                )
                console.log(response.data)
            },
        })
    })
