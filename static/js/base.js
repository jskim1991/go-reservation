'use strict'

const prompt = () => {
    const toast = (c) => {
        const { title = '', icon = 'success', position = 'top-end' } = c

        const Toast = Swal.mixin({
            toast: true,
            title,
            position: position,
            icon,
            showConfirmButton: false,
            timer: 3000,
            timerProgressBar: true,
            didOpen: (toast) => {
                toast.addEventListener('mouseenter', Swal.stopTimer)
                toast.addEventListener('mouseleave', Swal.resumeTimer)
            },
        })

        Toast.fire({})
    }

    const success = (c) => {
        const { title = '', text = '', confirmButtonText = '' } = c

        Swal.fire({
            title,
            text,
            icon: 'success',
            confirmButtonText,
        })
    }

    const error = (c) => {
        const { title = '', text = '', confirmButtonText = '' } = c

        Swal.fire({
            title,
            text,
            icon: 'error',
            confirmButtonText,
        })
    }

    const custom = async (c) => {
        const { icon = '', title = '', msg = '', showConfirmButton = true } = c
        const { value: result } = await Swal.fire({
            title,
            html: msg,
            icon,
            backdrop: false,
            focusConfirm: false,
            showCancelButton: true,
            showConfirmButton: showConfirmButton,
            willOpen: () => {
                if (c.willOpen !== undefined) {
                    c.willOpen()
                }
            },
            preConfirm: () => {
                return [
                    document.getElementById('start').value,
                    document.getElementById('end').value,
                ]
            },
            didOpen: () => {
                if (c.didOpen !== undefined) {
                    c.didOpen()
                }
            },
        })
    }

    const showDatepicker = async (c) => {
        const { msg = '', title = '' } = c
        const { value: result } = await Swal.fire({
            title,
            html: msg,
            backdrop: false,
            focusConfirm: false,
            showCancelButton: true,
            willOpen: () => {
                new DateRangePicker(
                    document.getElementById('reservation-dates-modal'),
                    {
                        format: 'yyyy-mm-dd',
                        showOnFocus: false,
                        autohide: true,
                        minDate: new Date(),
                    }
                )
            },
            didOpen: () => {
                document.getElementById('start').blur()
                document.querySelector('.swal2-actions').style.zIndex = '0' // prevent buttons getting in the way of datepicker
            },
            preConfirm: () => {
                return [
                    document.getElementById('start').value,
                    document.getElementById('end').value,
                ]
            },
        })

        if (result) {
            if (result.dismiss === Swal.DismissReason.cancel) {
                return
            }
            if (result.value === '') {
                return
            }
            if (c.callback === undefined) {
                return
            }

            c.callback(result)

            Swal.fire(JSON.stringify(result))
        }
    }

    return {
        toast,
        success,
        error,
        custom,
        showDatepicker,
    }
}

const validate = () => {
    const forms = document.querySelectorAll('.needs-validation')
    forms.forEach((form) => {
        form.addEventListener('submit', (e) => {
            if (form.checkValidity() == false) {
                e.preventDefault()
                e.stopPropagation()

                notification.error({
                    title: 'Required inputs missing',
                    text: 'Please enter values for all fields',
                    confirmButtonText: 'confirm',
                })
            }
            form.classList.add('was-validated')
        })
    })
}

const notify = (displayMessage, messageType) => {
    notie.alert({
        type: messageType,
        text: displayMessage,
        stay: false,
        time: 3,
        position: 'top',
    })
}

const showModal = (title, text, icon, confirmButtonText) => {
    Swal.fire({
        title,
        text,
        icon,
        confirmButtonText,
    })
}

const notification = prompt()
validate()
