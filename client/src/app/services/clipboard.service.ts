import { DOCUMENT } from '@angular/common';
import { HttpClient } from '@angular/common/http';
import { Inject, Injectable } from '@angular/core';
import { Observable, Subject } from 'rxjs';

interface ClipboardResponse {
  clipboard: string;
}

@Injectable({
  providedIn: 'root',
})
export class ClipboardService {
  private clipboardSubject = new Subject<string>();

  get clipboard(): Observable<string> {
    return this.clipboardSubject.asObservable();
  }

  private previousClipboardContents: string = "";

  constructor(private http: HttpClient) {
    http.get<ClipboardResponse>("api/clipboard").toPromise()
      .then((r) => {
        this.setNewClipboardContents(r.clipboard);
        this.startClipboardPoll();
      })
      .catch((e) => {
        console.error(e);
      })
  }

  private startClipboardPoll() {
    this.http.get<ClipboardResponse>("api/clipboardPoll").toPromise()
      .then((r) => {
        this.setNewClipboardContents(r.clipboard);
      })
      .catch((e) => {
        console.error(e);
      })
      .finally(() => {
        this.startClipboardPoll();
      })
  }

  private setNewClipboardContents(contents: string) {
    if (!contents || contents === this.previousClipboardContents) return;

    this.previousClipboardContents = contents;
    this.clipboardSubject.next(contents);
  }
}
