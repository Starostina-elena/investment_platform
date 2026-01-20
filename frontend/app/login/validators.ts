import {FormEvent} from "react";


export function LoginValidator(e: FormEvent<HTMLInputElement>) {
    const login = e.currentTarget.value;
    const match = login.match(/[a-zA-Z0-9а-яА-Я]+/);
    if (match && match[0].length !== login.length) {
        e.currentTarget.setCustomValidity("Логин содержит недопустимые символы (пробелы и любые спец.символы запрещены)");
        return;
    } else {
        e.currentTarget.setCustomValidity("");
    }
}

export function EmailValidator(e: FormEvent<HTMLInputElement>) {
    if (!e.currentTarget.value.match(/^(([^<>()\[\]\\.,;:\s@"]+(\.[^<>()\[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/)) {
        e.currentTarget.setCustomValidity("Неккоректный адрес электронной почты");
        e.currentTarget.setAttribute("aria-invalid", "true");
        return;
    } else {
        e.currentTarget.setCustomValidity("");
        e.currentTarget.setAttribute("aria-invalid", "false")
    }
}

export function PasswordValidator(e: FormEvent<HTMLInputElement>) {
    const password = e.currentTarget.value;
    const match = password.match(/^[a-zA-Z0-9а-яА-Я.&?$%*@#]+$/);
    if (!match) {
        const invertMatch = password.match(/[^a-zA-Z0-9а-яА-Я.&?$%*@#]/);
        if (invertMatch) {
            e.currentTarget.setCustomValidity("Пароль содержит недопустимый символ \"" + invertMatch[0] + "\", разрешены только латинские и русские буквы, цифры, спец.символы: . & ? $ % * @ #");
            e.currentTarget.setAttribute("aria-invalid", "true");
        }
        return;
    }
    e.currentTarget.setCustomValidity("");
    e.currentTarget.setAttribute("aria-invalid", "false");
}

export function NameValidator(e: FormEvent<HTMLInputElement>) {
    const name = e.currentTarget.value;
    const match = name.match(/^[a-zA-Zа-яА-Я]+$/);
    if (name.length > 0 && !match) {
        e.currentTarget.setCustomValidity("Поле содержит недопустимые символы (пробелы, цифры и любые спец.символы запрещены)");
        e.currentTarget.setAttribute("aria-invalid", "true")
        return;
    }
    e.currentTarget.setAttribute("aria-invalid", "false")
    e.currentTarget.setCustomValidity("");
}

export function NicknameValidator(e: FormEvent<HTMLInputElement>) {
    const nickname = e.currentTarget.value;
    const match = nickname.match(/^[a-zA-Zа-яА-Я0-9_()!@#$-]+$/);
    if (nickname.length > 0 && !match) {
        e.currentTarget.setCustomValidity("Поле содержит недопустимые символы (пробелы, спец.символы запрещены)");
        e.currentTarget.setAttribute("aria-invalid", "true")
        return;
    }
    e.currentTarget.setAttribute("aria-invalid", "false")
    e.currentTarget.setCustomValidity("");
}

export function isValidEmail(text: string) {
    return !(text.match(/^(([^<>()\[\]\\.,;:\s@"]+(\.[^<>()\[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/) === null)
}


export function isValidLogin(text: string) {
    const match = text.match(/[a-zA-Z0-9а-яА-Я]+/);
    if (match && match[0].length !== text.length) {
        return false
    } else {
        return true
    }
}