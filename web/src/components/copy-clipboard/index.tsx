/*
MIT License

Copyright (c) 2020 The KubeLens Authors

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/
import React from 'react';
import Clipboard from 'react-clipboard.js';
import ClipImg from '../../assets/copy-clipboard.png';
import './styles.css';

export type CopyClipboardProps = {
  size: number,
  value: string,
  labelText?: any,
  bottom?: number
};

const CopyClipboard = ({
  size,
  value,
  labelText,
  bottom
}: CopyClipboardProps) => {
  let style: any = { backgroundColor: 'transparent', border: 'none', cursor: 'copy' };
  if (bottom) {
    style.bottom = 0
  }

  const clip = (
    <Clipboard
      style={style}
      data-clipboard-text={value}>

      <img style={{ marginTop: '-5px' }} height={size} src={ClipImg} alt="Copy to Clipboard" />
    </Clipboard>
  );
  return (
    <div >
      {labelText ? labelText : null} {clip}
    </div>
  )
}

export default CopyClipboard;
