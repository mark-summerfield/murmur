" Vim syntax file
" Language:        Urm
" Author:          Mark Summerfield <mark@qtrac.eu>
" URL:             https://github.com/mark-summerfield/murmur
" Licence:         Public Domain
" Latest Revision: 2023-09-18

if exists("b:current_syntax")
  finish
endif

syn clear
syn sync fromstart linebreaks=3 minlines=50

" Order matters!

syn match urmComment /;.*/
syn match urmNumber /\d\+/
syn match urmLabel /[A-Za-z][A-Za-z0-9]*/
syn match urmSetLabel /^\S\+:/
syn match urmStartLabel /^START:/
syn match urmAddress /@[A-Za-z][A-Za-z0-9]*/
syn match urmCommand1 /[CDIJPSTZcdijpstz]\((\)\@=/
syn match urmCommand2 /STOP/

"" See https://sashamaps.net/docs/resources/20-colors/
hi urmComment	    guifg=forestgreen
hi urmAddress	    gui=italic guifg=#911EB4 "purple
hi urmLabel	    guifg=#9A6324 "brown
hi urmSetLabel	    guifg=#9A6324 "brown
hi urmStartLabel    gui=italic guifg=#469990 "teal
hi urmNumber	    guifg=#4363D8 "blue
hi urmCommand1 	    gui=bold guifg=darkblue
hi urmCommand2 	    gui=bold guifg=darkblue
