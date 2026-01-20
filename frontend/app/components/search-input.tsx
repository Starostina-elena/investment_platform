'use client'
import "./search-input.css";
import {useEffect, useRef, useState} from "react";

export function SearchInput({options, setValueAction, value, onChange, ...settings}:
                                {
                                    setValueAction: (key: string) => void,
                                    options: { [key: string]: string },
                                    value?: string
                                } & React.InputHTMLAttributes<HTMLInputElement>) {
    const [innerValue, setInnerValue] = useState('');
    const [focused, setFocused] = useState(false);
    const input = useRef<HTMLInputElement>(null);

    useEffect(() => {
        if (value && options[value])
            setInnerValue(value);
    }, [value]);

    return (
        <div className="chooser">
            <input type="text" style={focused || !options[innerValue] ? undefined : {color: 'transparent'}}
                   onFocus={() => {
                       setFocused(true);
                       if (innerValue) setInnerValue('');
                   }} ref={input}
                   value={innerValue} onBlur={() => setTimeout(() => setFocused(false), 100)}
                   onChange={e => {
                       setInnerValue(e.target.value);
                       if (!options[e.target.value])
                           e.target.setCustomValidity('Выберите элемент из списка!');
                       else
                           e.target.setCustomValidity('');
                       onChange && options[innerValue] && onChange(e)
                   }} {...settings}/>
            {!focused && options[innerValue] && <div className="chooser__label">
                {innerValue ? innerValue + ' - ' + options[innerValue] : ''}
            </div>}
            <div style={{position: 'relative', width: '100%'}}>
                <div style={{position: 'absolute', top: '0', left: 0}} className="chooser__list">
                    {focused && Object.keys(options).filter(e => e.startsWith(innerValue)).map(u =>
                        <div key={u}
                             onClick={e => {
                                 setValueAction(u);
                                 setInnerValue(u);
                                 e.currentTarget.blur();
                                 input.current?.setCustomValidity('');
                             }}>{u + ' - ' + options[u]}
                        </div>
                    )}
                </div>
            </div>
        </div>
    )
}