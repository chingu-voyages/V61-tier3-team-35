"use client"
import { useEffect, useState } from "react"
import Key from "./Key"
import KeyboardRows from "./keyboardRows"


export default function ({ activeKey, handleKeyPress, keyboardStatuses }) {

    return (
        <article className="flex flex-col items-center gap-1.5 py-10 w-full">
            {KeyboardRows.map((row, rowIndex) => (
                <div key={rowIndex} className="flex gap-1.5">
                    {row.map((key, index) => (
                        <Key className={activeKey === key ? "bg-gray-400/60" : ""} key={index} text={key} status={keyboardStatuses[key]} onClick={() => { handleKeyPress(key) }} />
                    ))}
                </div>
            ))}
        </article>
    )
}