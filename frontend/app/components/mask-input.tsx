import styles from './mask-input.module.css';
import {useEffect, useRef, useState} from 'react';

interface MaskInputProps {
    digits: number;
    mask?: string;
    value?: string;
    onChange?: React.ChangeEventHandler<HTMLInputElement>;
}

const ExtractDigits = (str: string) => Array.from(str.matchAll(/[0-9]/g));
const ApplyMask = (str: string, mask: string) => ((i) => mask
    .split('')
    .reduce((acc, e) => acc + (e == '_' && str[i] ? str[i++] : e), ''))(0);
const MaskInput = ({
                       digits,
                       mask = '',
                       value,
                       onChange,
                       ...rest
                   }: MaskInputProps & React.InputHTMLAttributes<HTMLInputElement>) => {
    const input = useRef<HTMLInputElement>(null);

    const masked = mask.split('').reduce((acc, e) => acc + +(e == '_'), 0);
    if (masked > digits) throw new Error('Mask is too long');
    const finalMask = mask + new Array(digits - masked).fill('_').join('');

    useEffect(() => {
        if (value && input.current && value != input.current.value) {
            input.current.value = value;
        }
    }, [value]);

    return (
        <input ref={input}
               type="text" data-old-value=''
               onChange={(e) => {
                   e.persist()
                   const caretStart = e.target.selectionStart;
                   const caretEnd = e.target.selectionEnd;

                   const value = ExtractDigits(e.currentTarget.value).slice(0, digits).join('');
                   const typedIncorrect = !!e.currentTarget.value.slice((caretStart || 0) - 1, caretStart || 0).match(/^[^0-9]+$/);

                   const maskedValue = ApplyMask(value, finalMask);
                   const isCarretTooForward = caretStart && caretStart > (e.currentTarget.getAttribute('data-old-value')?.indexOf('_') || -1);

                   e.currentTarget.setAttribute('data-old-value', ApplyMask(ExtractDigits(e.currentTarget.value).join(''), finalMask));
                   e.currentTarget.value = maskedValue;

                   if (onChange)
                       onChange(e);

                   if (value.length < digits && value.length > 0) {
                       e.currentTarget.setCustomValidity('Заполните поле!');
                       e.currentTarget.setAttribute('aria-invalid', 'true');
                   } else {
                       e.currentTarget.setCustomValidity('');
                       e.currentTarget.setAttribute('aria-invalid', 'false');
                   }

                   if (caretStart && typedIncorrect)
                       e.target.setSelectionRange(caretStart - 1, caretStart - 1);
                   else if (isCarretTooForward)
                       e.target.setSelectionRange(maskedValue.indexOf('_'), maskedValue.indexOf('_'));
                   else
                       e.target.setSelectionRange(caretStart, caretEnd);
               }}
               className={`${styles.input} ${rest.className || ''}`}
               {...rest}
        />
    );
};

export default MaskInput;