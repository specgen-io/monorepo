package genjava

import (
	"github.com/specgen-io/specgen/v2/sources"
	"strings"
)

func generateAdapters(thePackage Module) []sources.CodeFile {
	files := []sources.CodeFile{}
	files = append(files, *generateBigDecimalAdapter(thePackage))
	files = append(files, *generateLocalDateAdapter(thePackage))
	files = append(files, *generateLocalDateTimeAdapter(thePackage))
	files = append(files, *generateUuidAdapter(thePackage))
	return files
}

func generateBigDecimalAdapter(thePackage Module) *sources.CodeFile {
	code := `
package [[.PackageName]];

import com.squareup.moshi.*;
import okio.Okio;

import java.io.*;
import java.math.BigDecimal;
import java.nio.charset.StandardCharsets;

public class BigDecimalAdapter {
	@FromJson
	public BigDecimal fromJson(JsonReader reader) throws IOException {
		var token = reader.peek();
		if (token != JsonReader.Token.NUMBER) {
			throw new JsonDataException("BigDecimal should be represented as number in JSON, found: "+token.name());
		}
		var source = reader.nextSource();
		return new BigDecimal(new String(source.readByteArray(), StandardCharsets.UTF_8));
	}

	@ToJson
	public void toJson(JsonWriter writer, BigDecimal value) throws IOException {
		var source = Okio.source(new ByteArrayInputStream(value.toString().getBytes()));
		var buffer = Okio.buffer(source);
		writer.value(buffer);
	}
}
`

	code, _ = sources.ExecuteTemplate(code, struct{ PackageName string }{thePackage.PackageName})
	return &sources.CodeFile{
		Path:    thePackage.GetPath("BigDecimalAdapter.java"),
		Content: strings.TrimSpace(code),
	}
}

func generateLocalDateAdapter(thePackage Module) *sources.CodeFile {
	code := `
package [[.PackageName]];

import com.squareup.moshi.*;

import java.time.LocalDate;

public class LocalDateAdapter {
	@FromJson
	private LocalDate fromJson(String string) {
		return LocalDate.parse(string);
	}

	@ToJson
	private String toJson(LocalDate value) {
		return value.toString();
	}
}
`

	code, _ = sources.ExecuteTemplate(code, struct{ PackageName string }{thePackage.PackageName})
	return &sources.CodeFile{
		Path:    thePackage.GetPath("LocalDateAdapter.java"),
		Content: strings.TrimSpace(code),
	}
}

func generateLocalDateTimeAdapter(thePackage Module) *sources.CodeFile {
	code := `
package [[.PackageName]];

import com.squareup.moshi.*;

import java.time.LocalDateTime;

public class LocalDateTimeAdapter {
	@FromJson
	private LocalDateTime fromJson(String string) {
		return LocalDateTime.parse(string);
	}

	@ToJson
	private String toJson(LocalDateTime value) {
		return value.toString();
	}
}
`

	code, _ = sources.ExecuteTemplate(code, struct{ PackageName string }{thePackage.PackageName})
	return &sources.CodeFile{
		Path:    thePackage.GetPath("LocalDateTimeAdapter.java"),
		Content: strings.TrimSpace(code),
	}
}

func generateUuidAdapter(thePackage Module) *sources.CodeFile {
	code := `
package [[.PackageName]];

import com.squareup.moshi.*;

import java.util.UUID;

public class UuidAdapter {
	@FromJson
	private UUID fromJson(String string) {
		return UUID.fromString(string);
	}

	@ToJson
	private String toJson(UUID value) {
		return value.toString();
	}
}
`

	code, _ = sources.ExecuteTemplate(code, struct{ PackageName string }{thePackage.PackageName})
	return &sources.CodeFile{
		Path:    thePackage.GetPath("UuidAdapter.java"),
		Content: strings.TrimSpace(code),
	}
}
