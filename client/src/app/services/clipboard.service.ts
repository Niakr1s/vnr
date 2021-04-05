import { DOCUMENT } from '@angular/common';
import { Inject, Injectable } from '@angular/core';
import { Observable, Subject } from 'rxjs';

@Injectable({
  providedIn: 'root',
})
export class ClipboardService {
  private clipboardSubject = new Subject<string>();

  get clipboard(): Observable<string> {
    return this.clipboardSubject.asObservable();
  }

  private observer: MutationObserver;
  constructor(@Inject(DOCUMENT) document: Document) {
    this.observer = this.makeMutationObserver();

    const body = document.querySelector('body');
    if (!body) {
      throw new Error('body is null');
    }
    this.observer.observe(body, { childList: true });
  }

  private makeMutationObserver(): MutationObserver {
    return new MutationObserver(() => {
      const ps = document.getElementsByTagName('p');

      // removing all except last
      while (ps.length > 0) {
        const sentence = ps[0].textContent;
        ps[0].remove();

        if (!sentence) {
          continue;
        }

        this.clipboardSubject.next(sentence);
      }
    });
  }
}
